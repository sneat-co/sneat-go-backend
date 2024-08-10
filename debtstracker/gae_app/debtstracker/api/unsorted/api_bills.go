package unsorted

import (
	"context"
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus/dto"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/decimal"
	"net/http"
)

func HandleGetBill(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	billID := r.URL.Query().Get("id")
	if billID == "" {
		api.BadRequestError(c, w, errors.New("Missing id parameter"))
		return
	}
	bill, err := facade2debtus.GetBillByID(c, nil, billID)
	if err != nil {
		api.InternalError(c, w, err)
		return
	}
	billToResponse(c, w, authInfo.UserID, bill)
}

func HandleCreateBill(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	splitMode := models.SplitMode(r.PostFormValue("split"))
	if !models.IsValidBillSplit(splitMode) {
		api.BadRequestMessage(c, w, fmt.Sprintf("Split parameter has unkown value: %v", splitMode))
		return
	}
	amountStr := r.PostFormValue("amount")
	if amountStr == "" {
		api.BadRequestMessage(c, w, fmt.Sprintf("Missing required parameter: amount. %v", r.PostForm))
		return
	}
	amount, err := decimal.ParseDecimal64p2(amountStr)
	if err != nil {
		api.BadRequestError(c, w, err)
		return
	}
	var members []dto.BillMemberDto
	{
		membersJSON := r.PostFormValue("members")
		if err = ffjson.Unmarshal([]byte(membersJSON), &members); err != nil {
			api.BadRequestError(c, w, err)
			return
		}

	}
	if len(members) == 0 {
		api.BadRequestMessage(c, w, "No members has been provided")
		return
	}
	billEntity := models.NewBillEntity(models.BillCommon{
		Status:        models.BillStatusDraft,
		SplitMode:     splitMode,
		CreatorUserID: authInfo.UserID,
		Name:          r.PostFormValue("name"),
		Currency:      money.CurrencyCode(r.PostFormValue("currency")),
		AmountTotal:   amount,
	})

	var (
		totalByMembers decimal.Decimal64p2
	)

	contactIDs := make([]string, 0, len(members))
	memberUserIDs := make([]string, 0, len(members))

	for i, member := range members {
		if member.ContactID == "" && member.UserID == "" {
			api.BadRequestMessage(c, w, fmt.Sprintf("members[%d]: ContactID == 0 && UserID == 0", i))
			return
		}
		if member.ContactID != "" {
			contactIDs = append(contactIDs, member.ContactID)
		}
		if member.UserID != "" {
			memberUserIDs = append(memberUserIDs, member.UserID)
		}
	}

	var contacts []models.ContactEntry
	if len(contactIDs) > 0 {
		if contacts, err = facade2debtus.GetContactsByIDs(c, nil, contactIDs); err != nil {
			api.InternalError(c, w, err)
			return
		}
	}

	var memberUsers []*models.AppUser
	if len(memberUserIDs) > 0 {
		if memberUsers, err = facade2debtus.User.GetUsersByIDs(c, memberUserIDs); err != nil {
			api.InternalError(c, w, err)
			return
		}
	}

	billMembers := make([]models.BillMemberJson, len(members))
	for i, member := range members {
		if member.UserID != "" && member.ContactID != "" {
			api.BadRequestMessage(c, w, fmt.Sprintf("Member has both UserID and ContactID: %v, %v", member.UserID, member.ContactID))
			return
		}
		totalByMembers += member.Amount
		billMembers[i] = models.BillMemberJson{
			MemberJson: models.MemberJson{
				UserID: member.UserID,
				Shares: member.Share,
			},
			Percent:    member.Percent,
			Owes:       member.Amount,
			Adjustment: member.Adjustment,
		}
		if member.ContactID != "" {
			for _, contact := range contacts {
				if contact.ID == member.ContactID {
					contactName := contact.Data.FullName()
					billMembers[i].ContactByUser = models.MemberContactsJsonByUser{
						contact.Data.UserID: models.MemberContactJson{
							ContactID:   member.ContactID,
							ContactName: contactName,
						},
					}
					if billMembers[i].Name == "" {
						billMembers[i].Name = contactName
					}
					goto contactFound
				}
			}
			api.BadRequestError(c, w, fmt.Errorf("contact not found by ID=%v", member.ContactID))
			return
		contactFound:
		}
		if member.UserID != "" {
			for _, u := range memberUsers {
				if u.ID == member.UserID {
					billMembers[i].Name = u.Data.FullName()
					break
				}
			}
		}
	}
	if totalByMembers != amount {
		api.BadRequestMessage(c, w, fmt.Sprintf("Total amount is not equal to sum of member's amounts: %v != %v", amount, totalByMembers))
		return
	}

	billEntity.SplitMode = models.SplitModePercentage

	if err = billEntity.SetBillMembers(billMembers); err != nil {
		api.InternalError(c, w, err)
		return
	}

	var bill models.Bill
	err = facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
		bill, err = facade2debtus.Bill.CreateBill(c, tx, billEntity)
		return
	})

	if err != nil {
		api.InternalError(c, w, err)
		return
	}
	billToResponse(c, w, authInfo.UserID, bill)
}

func billToResponse(c context.Context, w http.ResponseWriter, userID string, bill models.Bill) {
	if userID == "" {
		api.InternalError(c, w, errors.New("Required parameter userID == 0."))
		return
	}
	if bill.ID == "" {
		api.InternalError(c, w, errors.New("Required parameter bill.ID is empty string."))
		return
	}
	if bill.Data == nil {
		api.InternalError(c, w, errors.New("Required parameter bill.BillEntity is nil."))
		return
	}
	billDto := dto.BillDto{
		ID:   bill.ID,
		Name: bill.Data.Name,
		Amount: money.Amount{
			Currency: money.CurrencyCode(bill.Data.Currency),
			Value:    decimal.Decimal64p2(bill.Data.AmountTotal),
		},
	}
	billMembers := bill.Data.GetBillMembers()
	members := make([]dto.BillMemberDto, len(billMembers))
	for i, billMember := range billMembers {
		members[i] = dto.BillMemberDto{
			UserID:     billMember.UserID,
			ContactID:  billMember.ContactByUser[userID].ContactID,
			Amount:     billMember.Owes,
			Adjustment: billMember.Adjustment,
			Share:      billMember.Shares,
		}
	}
	billDto.Members = members
	api.JsonToResponse(c, w, map[string]dto.BillDto{"Bill": billDto}) // TODO: Define DTO as need to clean BillMember.ContactByUser
}
