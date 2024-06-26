package unsorted

import (
	"encoding/json"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade/dto"
	"github.com/strongo/logus"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func getUserID(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) (userID string) {
	userID = authInfo.UserID

	if stringID := r.URL.Query().Get("user"); stringID != "" {
		if !authInfo.IsAdmin && userID != authInfo.UserID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	return
}

type UserCounterpartiesResponse struct {
	UserID         int64
	Counterparties []dto.ContactListDto
}

func HandleCreateCounterparty(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	name := strings.TrimSpace(r.PostForm.Get("name"))
	email := strings.TrimSpace(r.PostForm.Get("email"))
	tel := strings.TrimSpace(r.PostForm.Get("tel"))

	contactDetails := models.ContactDetails{
		Username: name,
	}
	if len(email) > 0 {
		contactDetails.EmailAddressOriginal = email
	}
	if len(tel) > 0 {
		telNumber, err := strconv.ParseInt(tel, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		contactDetails.PhoneNumber = telNumber
	}
	var err error
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		api.InternalError(c, w, err)
		return
	}
	var counterparty models.ContactEntry
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		counterparty, _, err = facade.CreateContact(c, tx, authInfo.UserID, contactDetails)
		return err
	})

	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
	_, _ = w.Write([]byte(counterparty.ID))
}

func getContactID(w http.ResponseWriter, query url.Values) string {
	counterpartyID := query.Get("id")
	if counterpartyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Required parameter 'id' is missing."))
	}
	return counterpartyID
}

func HandleGetContact(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	query := r.URL.Query()
	counterpartyID := getContactID(w, query)
	if counterpartyID == "" {
		return
	}
	counterparty, err := facade.GetContactByID(c, nil, counterpartyID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	contactToResponse(c, w, authInfo, counterparty)
}

func contactToResponse(c context.Context, w http.ResponseWriter, authInfo auth.AuthInfo, contact models.ContactEntry) {
	if !authInfo.IsAdmin && contact.Data.UserID != authInfo.UserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	transfers, hasMoreTransfers, err := dtdal.Transfer.LoadTransfersByContactID(c, contact.ID, 0, 100)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	counterpartyJson := dto.ContactDetailsDto{
		ContactListDto: dto.ContactListDto{
			Status: contact.Data.Status,
			ContactDto: dto.ContactDto{
				ID:     contact.ID,
				Name:   contact.Data.FullName(),
				UserID: contact.Data.UserID,
			},
		},
		TransfersResultDto: dto.TransfersResultDto{
			HasMoreTransfers: hasMoreTransfers,
			Transfers:        dto.TransfersToDto(authInfo.UserID, transfers),
		},
	}
	if contact.Data.BalanceJson != "" {
		balance := json.RawMessage(contact.Data.BalanceJson)
		counterpartyJson.Balance = &balance
	}
	if contact.Data.EmailAddressOriginal != "" {
		counterpartyJson.Email = &dto.EmailInfo{
			Address:     contact.Data.EmailAddressOriginal,
			IsConfirmed: contact.Data.EmailConfirmed,
		}
	}
	if contact.Data.PhoneNumber != 0 {
		counterpartyJson.Phone = &dto.PhoneInfo{
			Number:      contact.Data.PhoneNumber,
			IsConfirmed: contact.Data.PhoneNumberConfirmed,
		}
	}
	if len(contact.Data.GroupIDs) > 0 {
		for _, groupID := range contact.Data.GroupIDs {
			var group models.GroupEntry
			if group, err = dtdal.Group.GetGroupByID(c, nil, groupID); err != nil {
				api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
				return
			}
			for _, member := range group.Data.GetGroupMembers() {
				for _, memberContactID := range member.ContactIDs {
					if memberContactID == contact.ID {
						counterpartyJson.Groups = append(counterpartyJson.Groups, dto.ContactGroupDto{
							ID:           group.ID,
							Name:         group.Data.Name,
							MemberID:     memberContactID,
							MembersCount: group.Data.MembersCount,
						})
					}
				}
			}
		}
	}

	api.JsonToResponse(c, w, counterpartyJson)
}

//type CounterpartyTransfer struct {
//
//}

func HandleDeleteContact(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	logus.Debugf(c, "HandleDeleteContact()")
	//err := r.ParseForm()
	//if err != nil {
	//	BadRequestError(c, hashedWriter, err)
	//	return
	//}
	contactID := getContactID(w, r.URL.Query())
	if contactID == "" {
		return
	}
	logus.Debugf(c, "contactID: %v", contactID)
	if _, err := facade.DeleteContact(c, contactID); err != nil {
		api.InternalError(c, w, err)
		return
	}
	logus.Infof(c, "ContactEntry deleted: %v", contactID)
}

func HandleArchiveCounterparty(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	//err := r.ParseForm()
	//if err != nil {
	//	BadRequestError(c, hashedWriter, err)
	//	return
	//}
	contactID := getContactID(w, r.URL.Query())
	if contactID == "" {
		return
	}
	if contact, err := facade.ChangeContactStatus(c, contactID, models.STATUS_ARCHIVED); err != nil {
		api.InternalError(c, w, err)
		return
	} else {
		contactToResponse(c, w, authInfo, contact)
	}
}

func HandleActivateCounterparty(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	//err := r.ParseForm()
	//if err != nil {
	//	BadRequestError(c, hashedWriter, err)
	//	return
	//}

	contactID := getContactID(w, r.URL.Query())
	if contactID == "" {
		return
	}
	if contact, err := facade.ChangeContactStatus(c, contactID, models.STATUS_ACTIVE); err != nil {
		api.InternalError(c, w, err)
		return
	} else {
		contactToResponse(c, w, authInfo, contact)
	}
}

func HandleUpdateCounterparty(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	counterpartyID := getContactID(w, r.URL.Query())
	if counterpartyID == "" {
		return
	}
	values := make(map[string]string, len(r.PostForm))
	for k, vals := range r.PostForm {
		switch len(vals) {
		case 1:
			values[k] = vals[0]
		case 0:
			values[k] = vals[0]
		default:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("Too many values for '%v'.", k)))
			return
		}
	}

	if counterpartyEntity, err := facade.UpdateContact(c, counterpartyID, values); err != nil {
		api.InternalError(c, w, err)
		return
	} else {
		contactToResponse(c, w, authInfo, models.NewDebtusContact(counterpartyID, counterpartyEntity))
	}
}
