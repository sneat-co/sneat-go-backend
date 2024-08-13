package unsorted

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus/dto4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp/person"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type UserCounterpartiesResponse struct {
	UserID         int64
	Counterparties []dto4debtus.ContactListDto
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
	spaceID := r.URL.Query().Get("spaceID")

	contactDetails := dto4contactus.ContactDetails{
		NameFields: person.NameFields{
			UserName: name,
		},
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
	var debtusContact models4debtus.DebtusSpaceContactEntry
	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		_, _, debtusContact, err = facade4debtus.CreateContact(c, tx, authInfo.UserID, spaceID, contactDetails)
		return err
	})

	if err != nil {
		api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
	_, _ = w.Write([]byte(debtusContact.ID))
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
	contactID := getContactID(w, query)
	spaceID := query.Get("spaceID")
	if contactID == "" {
		return
	}

	var db dal.DB
	var err error
	if db, err = facade.GetDatabase(c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	contact := dal4contactus.NewContactEntry(spaceID, contactID)
	debtusContact := models4debtus.NewDebtusSpaceContactEntry(spaceID, contactID, nil)

	if err = db.GetMulti(c, []dal.Record{contact.Record, debtusContact.Record}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	contactToResponse(c, w, authInfo, contact, debtusContact)
}

func contactToResponse(
	ctx context.Context,
	w http.ResponseWriter,
	authInfo auth.AuthInfo,
	contact dal4contactus.ContactEntry,
	debtusContact models4debtus.DebtusSpaceContactEntry,
) {
	if !authInfo.IsAdmin && contact.Data.UserID != authInfo.UserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	transfers, hasMoreTransfers, err := dtdal.Transfer.LoadTransfersByContactID(ctx, contact.ID, 0, 100)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	counterpartyJson := dto4debtus.ContactDetailsDto{
		ContactListDto: dto4debtus.ContactListDto{
			Status: contact.Data.Status,
			ContactDto: dto4debtus.ContactDto{
				ID:     contact.ID,
				Name:   contact.Data.Names.GetFullName(),
				UserID: contact.Data.UserID,
			},
		},
		TransfersResultDto: dto4debtus.TransfersResultDto{
			HasMoreTransfers: hasMoreTransfers,
			Transfers:        dto4debtus.TransfersToDto(authInfo.UserID, transfers),
		},
	}
	if len(debtusContact.Data.Balance) > 0 {
		counterpartyJson.Balance = debtusContact.Data.Balance
	}

	//if contact.Data.EmailAddressOriginal != "" {
	//	counterpartyJson.Email = &dto4debtus.EmailInfo{
	//		Address:     contact.Data.EmailAddressOriginal,
	//		IsConfirmed: contact.Data.EmailConfirmed,
	//	}
	//}
	//if contact.Data.PhoneNumber != 0 {
	//	counterpartyJson.Phone = &dto4debtus.PhoneInfo{
	//		Number:      contact.Data.PhoneNumber,
	//		IsConfirmed: contact.Data.PhoneNumberConfirmed,
	//	}
	//}

	//if len(contact.Data.SpaceIDs) > 0 {
	//	err = errors.New("not implemented")
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//	for _, spaceID := range contact.Data.SpaceIDs {
	//		var group models4splitus.GroupEntry
	//		if group, err = dtdal.Group.GetGroupByID(ctx, nil, spaceID); err != nil {
	//			api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//			return
	//		}
	//		for _, member := range group.Data.GetGroupMembers() {
	//			for _, memberContactID := range member.ContactIDs {
	//				if memberContactID == contact.ContactID {
	//					counterpartyJson.Groups = append(counterpartyJson.Groups, dto4debtus.ContactGroupDto{
	//						ContactID:           group.ContactID,
	//						Name:         group.Data.Name,
	//						MemberID:     memberContactID,
	//						MembersCount: group.Data.MembersCount,
	//					})
	//				}
	//			}
	//		}
	//	}
	//}

	api4debtus.JsonToResponse(ctx, w, counterpartyJson)
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
	spaceID := r.URL.Query().Get("spaceID")
	if contactID == "" {
		return
	}
	logus.Debugf(c, "contactID: %v", contactID)
	userCtx := facade.NewUserContext("")
	if err := facade4debtus.DeleteContact(c, userCtx, spaceID, contactID); err != nil {
		api4debtus.InternalError(c, w, err)
		return
	}
	logus.Infof(c, "DebtusSpaceContactEntry deleted: %v", contactID)
}

func HandleArchiveCounterparty(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	//err := r.ParseForm()
	//if err != nil {
	//	BadRequestError(c, hashedWriter, err)
	//	return
	//}
	contactID := getContactID(w, r.URL.Query())
	spaceID := r.URL.Query().Get("spaceID")
	if contactID == "" {
		return
	}
	userCtx := facade.NewUserContext("")
	if contact, debtusContact, err := facade4debtus.ChangeContactStatus(c, userCtx, spaceID, contactID, const4debtus.StatusArchived); err != nil {
		api4debtus.InternalError(c, w, err)
		return
	} else {
		contactToResponse(c, w, authInfo, contact, debtusContact)
	}
}

func HandleActivateCounterparty(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	//err := r.ParseForm()
	//if err != nil {
	//	BadRequestError(c, hashedWriter, err)
	//	return
	//}

	contactID := getContactID(w, r.URL.Query())
	spaceID := r.URL.Query().Get("spaceID")
	userCtx := facade.NewUserContext("")
	if contactID == "" {
		return
	}
	if contact, debtusContact, err := facade4debtus.ChangeContactStatus(c, userCtx, spaceID, contactID, const4debtus.StatusActive); err != nil {
		api4debtus.InternalError(c, w, err)
		return
	} else {
		contactToResponse(c, w, authInfo, contact, debtusContact)
	}
}

func HandleUpdateCounterparty(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	counterpartyID := getContactID(w, r.URL.Query())
	if counterpartyID == "" {
		return
	}
	spaceID := r.URL.Query().Get("spaceID")
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

	if debtusContact, err := facade4debtus.UpdateContact(c, spaceID, counterpartyID, values); err != nil {
		api4debtus.InternalError(c, w, err)
		return
	} else {
		contact := dal4contactus.NewContactEntry(spaceID, debtusContact.ID)
		contactToResponse(c, w, authInfo, contact, debtusContact)
	}
}
