package unsorted

import (
	"context"
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus/dto4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/decimal"
	"net/http"
)

func HandleGetBill(c context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	billID := r.URL.Query().Get("id")
	if billID == "" {
		api4debtus.BadRequestError(c, w, errors.New("Missing id parameter"))
		return
	}
	bill, err := facade4splitus.GetBillByID(c, nil, billID)
	if err != nil {
		api4debtus.InternalError(c, w, err)
		return
	}
	billToResponse(c, w, authInfo.UserID, bill)
}

func HandleCreateBill(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	splitMode := models4splitus.SplitMode(r.PostFormValue("split"))
	if !models4splitus.IsValidBillSplit(splitMode) {
		api4debtus.BadRequestMessage(ctx, w, fmt.Sprintf("Split parameter has unkown value: %v", splitMode))
		return
	}
	spaceID := r.PostFormValue("spaceID")
	if spaceID == "" {
		api4debtus.BadRequestMessage(ctx, w, "Missing required parameter: spaceID")
		return
	}
	amountStr := r.PostFormValue("amount")
	if amountStr == "" {
		api4debtus.BadRequestMessage(ctx, w, fmt.Sprintf("Missing required parameter: amount. %v", r.PostForm))
		return
	}
	amount, err := decimal.ParseDecimal64p2(amountStr)
	if err != nil {
		api4debtus.BadRequestError(ctx, w, err)
		return
	}
	var members []dto4debtus.BillMemberDto
	{
		membersJSON := r.PostFormValue("members")
		if err = ffjson.Unmarshal([]byte(membersJSON), &members); err != nil {
			api4debtus.BadRequestError(ctx, w, err)
			return
		}

	}
	if len(members) == 0 {
		api4debtus.BadRequestMessage(ctx, w, "No members has been provided")
		return
	}
	billEntity := models4splitus.NewBillEntity(models4splitus.BillCommon{
		Status:        models4splitus.BillStatusDraft,
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
			api4debtus.BadRequestMessage(ctx, w, fmt.Sprintf("members[%d]: ContactID == 0 && UserID == 0", i))
			return
		}
		if member.ContactID != "" {
			contactIDs = append(contactIDs, member.ContactID)
		}
		if member.UserID != "" {
			memberUserIDs = append(memberUserIDs, member.UserID)
		}
	}

	var (
		debtusContacts []models4debtus.DebtusSpaceContactEntry
		contacts       []dal4contactus.ContactEntry
	)
	if len(contactIDs) > 0 {
		if debtusContacts, err = facade4debtus.GetDebtusSpaceContactsByIDs(ctx, nil, spaceID, contactIDs); err != nil {
			api4debtus.InternalError(ctx, w, err)
			return
		}
		if contacts, err = dal4contactus.GetContactsByIDs(ctx, nil, spaceID, contactIDs); err != nil {
			api4debtus.InternalError(ctx, w, err)
			return
		}
	}

	var memberUsers []dbo4userus.UserEntry
	if len(memberUserIDs) > 0 {
		if memberUsers, err = dal4userus.GetUsersByIDs(ctx, memberUserIDs); err != nil {
			api4debtus.InternalError(ctx, w, err)
			return
		}
	}

	billMembers := make([]*briefs4splitus.BillMemberBrief, len(members))
	for i, member := range members {
		if member.UserID != "" && member.ContactID != "" {
			api4debtus.BadRequestMessage(ctx, w, fmt.Sprintf("Member has both UserID and ContactID: %v, %v", member.UserID, member.ContactID))
			return
		}
		totalByMembers += member.Amount
		billMembers[i] = &briefs4splitus.BillMemberBrief{
			MemberBrief: briefs4splitus.MemberBrief{
				UserID: member.UserID,
				Shares: member.Share,
			},
			Percent:    member.Percent,
			Owes:       member.Amount,
			Adjustment: member.Adjustment,
		}
		if member.ContactID != "" {
			for contactIndex, debtusContact := range debtusContacts {
				if debtusContact.ID == member.ContactID {
					contactName := debtusContact.Data.FullName()
					contact := contacts[contactIndex]
					billMembers[i].ContactByUser = briefs4splitus.MemberContactBriefsByUserID{
						contact.Data.UserID: briefs4splitus.MemberContactBrief{
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
			api4debtus.BadRequestError(ctx, w, fmt.Errorf("debtusContact not found by ContactID=%v", member.ContactID))
			return
		contactFound:
		}
		if member.UserID != "" {
			for _, u := range memberUsers {
				if u.ID == member.UserID {
					billMembers[i].Name = u.Data.GetFullName()
					break
				}
			}
		}
	}
	if totalByMembers != amount {
		api4debtus.BadRequestMessage(ctx, w, fmt.Sprintf("Total amount is not equal to sum of member's amounts: %v != %v", amount, totalByMembers))
		return
	}

	billEntity.SplitMode = models4splitus.SplitModePercentage

	if err = billEntity.SetBillMembers(billMembers); err != nil {
		api4debtus.InternalError(ctx, w, err)
		return
	}

	var bill models4splitus.BillEntry
	err = facade.RunReadwriteTransaction(ctx, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
		bill, err = facade4splitus.CreateBill(ctx, tx, spaceID, billEntity)
		return
	})

	if err != nil {
		api4debtus.InternalError(ctx, w, err)
		return
	}
	billToResponse(ctx, w, authInfo.UserID, bill)
}

func billToResponse(c context.Context, w http.ResponseWriter, userID string, bill models4splitus.BillEntry) {
	if userID == "" {
		api4debtus.InternalError(c, w, errors.New("Required parameter userID == 0."))
		return
	}
	if bill.ID == "" {
		api4debtus.InternalError(c, w, errors.New("Required parameter bill.ContactID is empty string."))
		return
	}
	if bill.Data == nil {
		api4debtus.InternalError(c, w, errors.New("Required parameter bill.BillDbo is nil."))
		return
	}
	billDto := dto4debtus.BillDto{
		ID:   bill.ID,
		Name: bill.Data.Name,
		Amount: money.Amount{
			Currency: money.CurrencyCode(bill.Data.Currency),
			Value:    decimal.Decimal64p2(bill.Data.AmountTotal),
		},
	}
	billMembers := bill.Data.GetBillMembers()
	members := make([]dto4debtus.BillMemberDto, len(billMembers))
	for i, billMember := range billMembers {
		members[i] = dto4debtus.BillMemberDto{
			UserID:     billMember.UserID,
			ContactID:  billMember.ContactByUser[userID].ContactID,
			Amount:     billMember.Owes,
			Adjustment: billMember.Adjustment,
			Share:      billMember.Shares,
		}
	}
	billDto.Members = members
	api4debtus.JsonToResponse(c, w, map[string]dto4debtus.BillDto{"BillEntry": billDto}) // TODO: Define DTO as need to clean BillMember.ContactByUser
}
