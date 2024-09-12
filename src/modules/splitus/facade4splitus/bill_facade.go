package facade4splitus

import (
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/facade4anybot"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"math"

	"strconv"
	"time"

	"context"
	"github.com/strongo/decimal"
)

func AssignBillToGroup(ctx context.Context, tx dal.ReadwriteTransaction, inBill models4splitus.BillEntry, spaceID, userID string) (bill models4splitus.BillEntry, splitusSpace models4splitus.SplitusSpaceEntry, err error) {
	bill = inBill
	if err = bill.Data.AssignToGroup(spaceID); err != nil {
		return
	}
	splitusSpace = models4splitus.NewSplitusSpaceEntry(spaceID)
	if len(bill.Data.Members) == 0 {
		if err = tx.Get(ctx, splitusSpace.Record); err != nil {
			return
		}
		if len(splitusSpace.Data.Members) > 0 {
			groupMembers := splitusSpace.Data.GetGroupMembers()

			billMembers := make([]*briefs4splitus.BillMemberBrief, len(groupMembers))
			paidIsSet := false
			for i, gm := range groupMembers {
				billMembers[i] = &briefs4splitus.BillMemberBrief{
					MemberBrief: gm.MemberBrief,
				}
				billMembers[i].AddedByUserID = userID
				if gm.UserID == bill.Data.CreatorUserID {
					billMembers[i].Paid = bill.Data.AmountTotal
					paidIsSet = true
				}
			}
			if !paidIsSet {
				for i, bm := range billMembers {
					if bm.UserID == userID {
						billMembers[i].Paid = bill.Data.AmountTotal
						paidIsSet = true
						break
					}
				}
				if !paidIsSet { // current user is not members of the bill
					//splitusSpace.AddOrGetMember(userID, "", )
					user := dbo4userus.NewUserEntry(userID)
					if err = dal4userus.GetUser(ctx, tx, user); err != nil {
						return
					}
					_, _, _, groupMember, _ := splitusSpace.Data.AddOrGetMember(userID, "", user.Data.GetFullName())

					billMembers = append(billMembers, &briefs4splitus.BillMemberBrief{
						MemberBrief: groupMember.MemberBrief,
						Paid:        bill.Data.AmountTotal,
					})
				}
			}
			if err = bill.Data.SetBillMembers(billMembers); err != nil {
				return
			}
			if bill.Data.Currency != money.CurrencyCode("") {
				if _, err = splitusSpace.Data.ApplyBillBalanceDifference(bill.Data.Currency, bill.Data.GetBalance().BillBalanceDifference(briefs4splitus.BillBalanceByMember{})); err != nil {
					return
				}
				if err = tx.Set(ctx, splitusSpace.Record); err != nil {
					return
				}
			}
			logus.Debugf(ctx, "bill.GetBillMembers(): %+v", bill.Data.GetBillMembers())
		}
	}
	return
}

func CreateBill(ctx context.Context, tx dal.ReadwriteTransaction, spaceID string, billEntity *models4splitus.BillDbo) (bill models4splitus.BillEntry, err error) {
	if ctx == nil {
		panic("Parameter ctx context.Context is required")
	}
	logus.Debugf(ctx, "billFacade.CreateBill(%v)", billEntity)
	if tx == nil {
		panic("parameter tx dal.ReadwriteTransaction is required")
	}
	if billEntity == nil {
		panic("Parameter billEntity *models.BillDbo is required")
	}
	if !models4splitus.IsValidBillSplit(billEntity.SplitMode) {
		panic(fmt.Sprintf("billEntity.SplitMode has unknown value: %v", billEntity.SplitMode))
	}
	if billEntity.CreatorUserID == "" {
		err = fmt.Errorf("%w: billEntity.CreatorUserID == 0", facade4anybot.ErrBadInput)
		return
	}
	if billEntity.SplitMode == "" {
		err = fmt.Errorf("%w: Missing required property SplitMode", facade4anybot.ErrBadInput)
		return
	}
	if billEntity.AmountTotal == 0 {
		err = fmt.Errorf("%w: billEntity.AmountTotal == 0", facade4anybot.ErrBadInput)
		return
	}
	if billEntity.AmountTotal < 0 {
		err = fmt.Errorf("%w: billEntity.AmountTotal < 0: %v", facade4anybot.ErrBadInput, billEntity.AmountTotal)
		return
	}
	if billEntity.Status == "" {
		err = fmt.Errorf("%w: billEntity.Status property is required", err)
		return
	}
	if !models4splitus.IsValidBillStatus(billEntity.Status) {
		err = fmt.Errorf("%w: invalid status: %v, expected one of %v", facade4anybot.ErrBadInput, billEntity.Status, models4splitus.BillStatuses)
		return
	}

	billEntity.CreatedAt = time.Now()

	members := billEntity.GetBillMembers()
	//if len(members) == 0 {
	//	return bill, fmt.Errorf("len(members) == 0, MembersJson: %v", billEntity.MembersJson)
	//}

	if len(members) == 0 {
		billEntity.SplitMode = models4splitus.SplitModeEqually
	} else {
		contactIDs := make([]string, 0, len(members)-1)

		var (
			totalPercentageByMembers decimal.Decimal64p2
			totalSharesPerMembers    int
			totalPaidByMembers       decimal.Decimal64p2
			totalOwedByMembers       decimal.Decimal64p2
			payersCount              int
			equalAmount              decimal.Decimal64p2
			shareAmount              decimal.Decimal64p2
		)

		switch billEntity.SplitMode {
		case models4splitus.SplitModeShare:
			shareAmount = decimal.NewDecimal64p2FromFloat64(
				math.Floor(billEntity.AmountTotal.AsFloat64()/float64(len(members))*100+0.5) / 100,
			)
		case models4splitus.SplitModeEqually:
			amountToSplitEqually := billEntity.AmountTotal
			var totalAdjustmentByMembers decimal.Decimal64p2
			for i, member := range members {
				if member.Adjustment > billEntity.AmountTotal {
					err = fmt.Errorf("%w: members[%d].Adjustment > billEntity.AmountTotal", facade4anybot.ErrBadInput, i)
					return
				} else if member.Adjustment < 0 && member.Adjustment < -1*billEntity.AmountTotal {
					err = fmt.Errorf("%w: members[%d].AdjustmentInCents < 0 && AdjustmentInCents < -1*billEntity.AmountTotal", facade4anybot.ErrBadInput, i)
					return
				}
				totalAdjustmentByMembers += member.Adjustment
			}
			if totalAdjustmentByMembers > billEntity.AmountTotal {
				err = errors.New("totalAdjustmentByMembers > billEntity.AmountTotal")
				return
			}
			amountToSplitEqually -= totalAdjustmentByMembers
			equalAmount = decimal.NewDecimal64p2FromFloat64(
				math.Floor(amountToSplitEqually.AsFloat64()/float64(len(members))*100+0.5) / 100,
			)
		}

		// We use it to check equal split
		amountsCountByValue := make(map[decimal.Decimal64p2]int)

		// Calculate totals & initial checks
		for i, member := range members {
			if member.Paid != 0 {
				payersCount += 1
				totalPaidByMembers += member.Paid
			}
			totalOwedByMembers += member.Owes
			totalSharesPerMembers += member.Shares

			// Individual member checks - we can't move this checks down as it should fail first before deviation checks
			{
				if member.Owes < 0 {
					err = fmt.Errorf("%w: members[%d].Owes is negative: %v", facade4anybot.ErrBadInput, i, member.Owes)
					return
				}
				if member.UserID != billEntity.CreatorUserID {
					if len(member.ContactByUser) == 0 {
						err = fmt.Errorf("len(members[i].ContactByUser) == 0: i==%v", i)
						return
					}
					if member.UserID == "" {
						if len(member.ContactByUser) == 0 {
							err = errors.New("bill member is missing ContactByUser ContactID")
							return
						}

						for _, counterparty := range member.ContactByUser {
							if counterparty.ContactID == "" {
								panic("counterparty.ContactID == 0")
							}
							if err != nil {
								return
							}
							var duplicateContactID bool
							for _, cID := range contactIDs {
								if cID == counterparty.ContactID {
									duplicateContactID = true
									break
								}
							}
							if !duplicateContactID {
								contactIDs = append(contactIDs, counterparty.ContactID)
							}
						}
					}
				}
			}
		}

		adjustmentsCount := 0
		for i, member := range members {
			if member.Adjustment != 0 {
				adjustmentsCount++
			}
			ensureNoAdjustment := func() {
				if member.Adjustment != 0 {
					panic(fmt.Sprintf("Member #%d has Adjustment property not allowed with split mode %v", i, billEntity.SplitMode))
				}
			}
			ensureNoShare := func() {
				if member.Shares != 0 {
					panic(fmt.Sprintf("Member #%d has Shares property not allowed with split mode %v", i, billEntity.SplitMode))
				}
			}
			ensureEqualShare := func() {
				if member.Shares != members[0].Shares {
					panic(fmt.Sprintf("members[%d] has Shares not equal to members[0].Shares: %d != %d", i, member.Shares, members[i].Shares))
				}
			}

			ensureMemberAmountDeviateWithin1cent := func() error {
				//if totalOwedByMembers == 0 && totalOwedByMembers == 0 {
				//	return nil
				//}
				switch billEntity.SplitMode {
				case models4splitus.SplitModeShare:
					expectedAmount := int64(shareAmount) * int64(member.Shares)
					deviation := expectedAmount - int64(member.Owes)
					if deviation > 1 || deviation < -1 {
						return fmt.Errorf("%w: member #%d has amount %v deviated too much (for %v) from expected %v", facade4anybot.ErrBadInput, i, member.Owes, decimal.Decimal64p2(deviation), decimal.Decimal64p2(expectedAmount))
					}
				default:
					deviation := int64(member.Owes - member.Adjustment - equalAmount)
					if deviation > 1 || deviation < -1 {
						return fmt.Errorf("%w: member #%d has amount %v deviated too much (for %v) from equal %v", facade4anybot.ErrBadInput, i, member.Owes, decimal.Decimal64p2(deviation), equalAmount)
					}
				}
				return nil
			}
			switch billEntity.SplitMode {
			case models4splitus.SplitModeEqually:
				// ensureNoAdjustment()
				ensureEqualShare()
				if err = ensureMemberAmountDeviateWithin1cent(); err != nil {
					return
				}
				amountsCountByValue[member.Owes]++
			case models4splitus.SplitModeExactAmount:
				ensureNoAdjustment()
				ensureNoShare()
			case models4splitus.SplitModePercentage:
				totalPercentageByMembers += member.Percent
				// ensureNoAdjustment()
			case models4splitus.SplitModeShare:
				if member.Shares == 0 {
					err = fmt.Errorf("%w: member %d is missing Shares value", facade4anybot.ErrBadInput, i)
					return
				}
				// ensureNoAdjustment()
			}
		}

		if payersCount > 1 {
			err = ErrBillHasTooManyPayers
			return
		}

		if !(billEntity.Status == const4debtus.StatusDraft && totalPaidByMembers == 0) && totalPaidByMembers != billEntity.AmountTotal {
			err = fmt.Errorf("%w: total paid for all members should be equal to billEntity amount (%d), got %d", facade4anybot.ErrBadInput, billEntity.AmountTotal, totalPaidByMembers)
			return
		}
		switch billEntity.SplitMode {
		case models4splitus.SplitModeEqually:
			if len(amountsCountByValue) > 2+adjustmentsCount {
				err = fmt.Errorf("%w: len(amountsCountByValue):%d > 2 + adjustmentsCount:%d", facade4anybot.ErrBadInput, amountsCountByValue, adjustmentsCount)
				return
			}
		case models4splitus.SplitModePercentage:
			if totalPercentageByMembers != decimal.NewDecimal64p2FromInt(100) {
				err = fmt.Errorf("%w: total percentage for all members should be 100%%, got %v%%", facade4anybot.ErrBadInput, totalPercentageByMembers)
				return
			}
		case models4splitus.SplitModeShare:
			if billEntity.Shares == 0 {
				billEntity.Shares = totalSharesPerMembers
			} else if billEntity.Shares != totalSharesPerMembers {
				err = fmt.Errorf("%w: billEntity.Shares != totalSharesPerMembers", facade4anybot.ErrBadInput)
				return
			}
		}

		if (totalOwedByMembers != 0 || totalPaidByMembers != 0) && totalOwedByMembers != billEntity.AmountTotal {
			err = fmt.Errorf("totalOwedByMembers != billEntity.AmountTotal: %v != %v", totalOwedByMembers, billEntity.AmountTotal)
			return
		}

		// Load counterparties so we can get respective userIDs
		var counterparties []models4debtus.DebtusSpaceContactEntry
		// Use non transactional context
		counterparties, err = facade4debtus.GetDebtusSpaceContactsByIDs(ctx, tx, spaceID, contactIDs)
		if err != nil {
			err = fmt.Errorf("failed to get counterparties by ContactID: %w", err)
			return
		}

		// Assign userIDs from counterparty to respective member
		for _, member := range members {
			for _, counterparty := range counterparties {
				// TODO: assign not just for creator?
				if member.UserID == "" && member.ContactByUser[billEntity.CreatorUserID].ContactID == counterparty.ID {
					member.UserID = counterparty.Data.CounterpartyUserID
					break
				}
			}
		}

		billEntity.ContactIDs = contactIDs[:]
	}

	if bill, err = InsertBillEntity(ctx, tx, billEntity); err != nil {
		return
	}

	billHistory := models4splitus.NewBillHistoryBillCreated(bill, nil)
	if err = dtdal.InsertWithRandomStringID(ctx, tx, billHistory.Record); err != nil {
		return
	}
	return
}

//func (billFacade) CreateBillTransfers(ctx context.Context, billID string) error {
//	bill, err := facade4debtus.GetBillByID(ctx, billID)
//	if err != nil {
//		return err
//	}
//
//	members := bill.GetBillMembers()
//
//	{ // Verify payers count
//		payersCount := 0
//		for _, member := range members {
//			if member.Paid != 0 {
//				payersCount += 1
//			}
//		}
//		if payersCount == 0 {
//			return ErrBillHasNoPayer
//		} else if payersCount > 1 {
//			return ErrBillHasTooManyPayers
//		}
//	}
//
//	for _, member := range members {
//		if member.Paid == 0 {
//			creatorContactID := member.ContactByUser[bill.CreatorUserID].ContactID
//			if err = BillEntry.createBillTransfer(c, billID, strconv.FormatInt(bill.CreatorUserID, 10)); err != nil {
//				return errors.Wrapf(err, "Failed to create bill trasfer for %d", creatorContactID)
//			}
//		}
//	}
//	return nil
//}
//
//func (billFacade) createBillTransfer(ctx context.Context, billID string, creatorCounterpartyID int64) error {
//	err := dtdal.DB.RunInTransaction(ctx, func(ctx context.Context) error {
//		bill, err := facade4debtus.GetBillByID(ctx, billID)
//
//		if err != nil {
//			return err
//		}
//		members := bill.GetBillMembers()
//
//		var (
//			borrower *models.BillMemberBrief
//			payer    *models.BillMemberBrief
//		)
//		sCreatorUserID := strconv.FormatInt(bill.CreatorUserID, 10)
//		for _, member := range members {
//			if member.Paid > 0 {
//				if payer != nil {
//					return ErrBillHasTooManyPayers
//				}
//				payer = &member
//				if borrower != nil {
//					break
//				}
//			} else if member.ContactByUser[sCreatorUserID].ContactID == creatorCounterpartyID {
//				borrower = &member
//				if payer != nil {
//					break
//				}
//			}
//		}
//		if borrower == nil {
//			return errors.New("BillEntry member not found by creatorCounterpartyID")
//		}
//		if payer == nil {
//			return ErrBillHasNoPayer
//		}
//		//transferSource := dtdal.NewTransferSourceBot("api4debtus", "no-id", "0") // TODO: Needs refactoring! Move it out of DAL, do we really need an interface?
//
//		from := models.TransferCounterpartyInfo{
//			UserID:    payer.UserID,
//			ContactID: payer.ContactByUser[sCreatorUserID].ContactID,
//		}
//		to := models.TransferCounterpartyInfo{
//			UserID:    borrower.UserID,
//			ContactID: payer.ContactByUser[sCreatorUserID].ContactID,
//		}
//		logus.Debugf(ctx, "from: %v", from)
//		logus.Debugf(ctx, "to: %v", to)
//		//_, _, _, _, _, _, err = CreateTransfer(
//		//	c,
//		//	strongoapp.EnvUnknown,
//		//	transferSource,
//		//	bill.CreatorUserID,
//		//	billID,
//		//	false,
//		//	0,
//		//	from, to,
//		//	 money.AmountTotal{Currency: money.CurrencyCode(bill.Currency), Value: bill.AmountTotal},
//		//	time.Time{},
//		//)
//		//if err != nil {
//		//	return err
//		//}
//		return nil
//	}, dtdal.CrossGroupTransaction)
//	return err
//}

type BillMemberUserInfo struct {
	ContactID string
	Name      string
}

func GetBillMembersUserInfo(_ context.Context, bill models4splitus.BillEntry, forUserID int64) (billMembersUserInfo []BillMemberUserInfo, err error) {
	sUserID := strconv.FormatInt(forUserID, 10)

	for i, member := range bill.Data.GetBillMembers() {
		var (
			billMemberContactJson briefs4splitus.MemberContactBrief
			ok                    bool
		)
		if billMemberContactJson, ok = member.ContactByUser[sUserID]; !ok {
			err = fmt.Errorf("Member  #%d does not have information for %v", i, sUserID)
			return
		}
		billMembersUserInfo[i] = BillMemberUserInfo{
			ContactID: billMemberContactJson.ContactID,
			Name:      billMemberContactJson.ContactName,
		}
	}
	return
}

func AddBillMember(
	ctx context.Context, tx dal.ReadwriteTransaction, userID string, inBill models4splitus.BillEntry, memberID, memberUserID string, memberUserName string, paid decimal.Decimal64p2,
) (
	bill models4splitus.BillEntry, splitusSpace models4splitus.SplitusSpaceEntry, changed, isJoined bool, err error,
) {
	if tx == nil {
		panic("This method should be called within transaction")
	}
	logus.Debugf(ctx, "billFacade.AddBillMember(bill.ContactID=%v, memberID=%v, memberUserID=%v, memberUserName=%v, paid=%v)", bill.ID, memberID, memberUserID, memberUserName, paid)
	if paid < 0 {
		panic("paid < 0")
	}
	bill = inBill
	if bill.ID == "" {
		panic("bill.ContactID is empty string")
	}

	// TODO: Verify bill was obtained within transaction

	previousBalance := bill.Data.GetBalance()

	var (
		//isNew bool
		index              int
		groupChanged       bool
		groupMember        briefs4splitus.SpaceSplitMember
		billMember         *briefs4splitus.BillMemberBrief
		billMembers        []*briefs4splitus.BillMemberBrief
		groupMembers       []briefs4splitus.SpaceSplitMember
		splitMembersBefore []briefs4splitus.SpaceSplitMember
	)

	totalAboutBefore := bill.Data.AmountTotal

	splitusSpace = models4splitus.NewSplitusSpaceEntry(bill.Data.SpaceID)
	if err = tx.Get(ctx, splitusSpace.Record); err != nil {
		return
	}
	copy(splitusSpace.Data.Members, splitMembersBefore)

	if _, groupChanged, _, groupMember, groupMembers = splitusSpace.Data.AddOrGetMember(memberUserID, "", memberUserName); groupChanged {
		splitusSpace.Data.SetGroupMembers(groupMembers)
	} else {
		logus.Debugf(ctx, "GroupEntry billMembers not changed, groupMember.ContactID: "+groupMember.ID)
	}

	_, changed, index, billMember, billMembers = bill.Data.AddOrGetMember(groupMember.ID, memberUserID, "", memberUserName)

	logus.Debugf(ctx, "billMember.ContactID: "+billMember.ID)

	if paid > 0 {
		if billMember.Paid == paid {
			// Already set
		} else if paid == bill.Data.AmountTotal {
			for i := range billMembers {
				billMembers[i].Paid = 0
			}
			billMember.Paid = paid
			changed = true
		} else {
			paidTotal := paid
			for _, bm := range billMembers {
				paidTotal += bm.Paid
			}
			if paidTotal <= bill.Data.AmountTotal {
				billMember.Paid = paid
				changed = true
			} else {
				err = errors.New("Total paid by members exceeds bill amount")
				return
			}
		}
	}
	if !changed {
		return
	}

	billMembers[index] = billMember

	logus.Debugf(ctx, "billMembers: %+v", billMembers)
	if err = bill.Data.SetBillMembers(billMembers); err != nil {
		return
	}
	logus.Debugf(ctx, "bill.GetBillMembers(): %+v", bill.Data.GetBillMembers())

	if err = SaveBill(ctx, tx, bill); err != nil {
		return
	}

	logus.Debugf(ctx, "bill.GetBillMembers() after save: %v", bill.Data.GetBillMembers())

	currentBalance := bill.Data.GetBalance()

	if balanceDifference := currentBalance.BillBalanceDifference(previousBalance); balanceDifference.IsNoDifference() {
		logus.Debugf(ctx, "BillEntry balanceDifference: %v", balanceDifference)
		if groupChanged, err = splitusSpace.Data.ApplyBillBalanceDifference(bill.Data.Currency, balanceDifference); err != nil {
			err = fmt.Errorf("failed to apply bill difference: %w", err)
			return
		}
		if groupChanged {
			if err = tx.Set(ctx, splitusSpace.Record); err != nil {
				return
			}
		}
	}

	logus.Debugf(ctx, "splitusSpace: %+v", splitusSpace)
	var splitMembersAfter []briefs4splitus.SpaceSplitMember
	copy(splitMembersAfter, splitusSpace.Data.Members)
	billHistory := models4splitus.NewBillHistoryMemberAdded(userID, bill, totalAboutBefore, splitMembersBefore, splitMembersAfter)

	if err = tx.Insert(ctx, billHistory.Record); err != nil {
		return
	}

	isJoined = true
	return
}

var (
	ErrSettledBillsCanNotBeDeleted   = errors.New("settled bills can't be deleted")
	ErrOnlyDeletedBillsCanBeRestored = errors.New("only deleted bills can be restored")
)

func DeleteBill(ctx context.Context, billID string, userID string) (bill models4splitus.BillEntry, err error) {
	if err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if bill, err = GetBillByID(ctx, nil, billID); err != nil {
			return
		}
		if bill.Data.Status == models4splitus.BillStatusSettled {
			err = ErrSettledBillsCanNotBeDeleted
			return
		}
		if bill.Data.Status == models4splitus.BillStatusDraft || bill.Data.Status == models4splitus.BillStatusOutstanding {
			billHistory := models4splitus.NewBillHistoryBillDeleted(userID, bill)
			if err = tx.Insert(ctx, billHistory.Record); err != nil {
				return
			}
			bill.Data.Status = models4splitus.BillStatusDeleted
			if err = SaveBill(ctx, tx, bill); err != nil {
				return
			}
		}
		if spaceID := bill.Data.GetUserGroupID(); spaceID != "" {
			splitusSpace := models4splitus.NewSplitusSpaceEntry(spaceID)
			if err = tx.Get(ctx, splitusSpace.Record); err != nil {
				return err
			}
			outstandingBills := splitusSpace.Data.GetOutstandingBills()
			delete(outstandingBills, billID)
			panic("not implemented yet")
			//for briefBillID, billJson := range outstandingBills {
			//	if briefBillID == billID {
			//		outstandingBills = append(outstandingBills[:i], outstandingBills[i+1:]...)
			//		if err = splitusSpace.Data.SetOutstandingBills(outstandingBills); err != nil {
			//			return err
			//		}
			//		groupMembers := splitusSpace.Data.GetGroupMembers()
			//		billMembers := bill.Data.GetBillMembers()
			//		for j, groupMember := range groupMembers {
			//			for _, billMember := range billMembers {
			//				if billMember.ContactID == groupMember.ContactID {
			//					groupMember.Balance[bill.Data.Currency] -= billMember.Balance()
			//					groupMembers[j] = groupMember
			//					break
			//				}
			//			}
			//		}
			//		splitusSpace.Data.SetGroupMembers(groupMembers)
			//		if err = dtdal.Group.SaveGroup(ctx, tx, splitusSpace); err != nil {
			//			return
			//		}
			//		break
			//	}
			//}
		}
		return
	}, dal.TxWithCrossGroup()); err != nil {
		return
	}
	return
}

func RestoreBill(ctx context.Context, billID string, userID string) (bill models4splitus.BillEntry, err error) {
	if err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if bill, err = GetBillByID(ctx, nil, billID); err != nil {
			return
		}
		if bill.Data.Status != models4splitus.BillStatusDeleted {
			err = ErrOnlyDeletedBillsCanBeRestored
			return
		}

		if len(bill.Data.Members) > 1 {
			bill.Data.Status = models4splitus.BillStatusOutstanding
		} else {
			bill.Data.Status = models4splitus.BillStatusDraft
		}
		billHistory := models4splitus.NewBillHistoryBillRestored(userID, bill)
		if err = tx.Insert(ctx, billHistory.Record); err != nil {
			return
		}
		if err = SaveBill(ctx, tx, bill); err != nil {
			return
		}
		if spaceID := bill.Data.SpaceID; spaceID != "" {
			splitusSpace := models4splitus.NewSplitusSpaceEntry(spaceID)
			if err = tx.Get(ctx, splitusSpace.Record); err != nil {
				return
			}
			var groupChanged bool
			if groupChanged, err = splitusSpace.Data.AddBill(bill); err != nil {
				return
			} else if groupChanged {
				if err = tx.Set(ctx, splitusSpace.Record); err != nil {
					return
				}
			}
		}
		return
	}, dal.TxWithCrossGroup()); err != nil {
		return
	}
	return
}

func GetBillByID(ctx context.Context, tx dal.ReadSession, billID string) (bill models4splitus.BillEntry, err error) {
	if tx == nil {
		if tx, err = facade.GetSneatDB(ctx); err != nil {
			return bill, err
		}
	}
	bill.ID = billID
	bill.Key = dal.NewKeyWithID(models4splitus.BillKind, billID)
	bill.Data = new(models4splitus.BillDbo)
	bill.Record = dal.NewRecordWithData(bill.Key, bill.Data)
	err = tx.Get(ctx, bill.Record)
	return
}

func InsertBillEntity(ctx context.Context, tx dal.ReadwriteTransaction, billEntity *models4splitus.BillDbo) (bill models4splitus.BillEntry, err error) {
	if billEntity == nil {
		panic("billEntity == nil")
	}
	if billEntity.CreatorUserID == "" {
		panic("CreatorUserID == 0")
	}
	if billEntity.AmountTotal == 0 {
		panic("AmountTotal == 0")
	}

	billEntity.CreatedAt = time.Now()
	bill.Data = billEntity

	err = tx.Insert(ctx, bill.Record)
	return
}

//func (billFacade billFacade) createTransfers(ctx context.Context, splitID int64) error {
//	split, err := dtdal.Split.GetSplitByID(ctx, splitID)
//	if err != nil {
//		return err
//	}
//	bills, err := dtdal.BillEntry.GetBillsByIDs(ctx, split.BillIDs)
//
//	balances := billFacade.getBalances(splitID, bills)
//	balances = billFacade.cleanupBalances(balances)
//
//	for currency, totalsByMember := range balances {
//		for memberID, memberTotal := range totalsByMember {
//			if memberTotal.Balance() > 0 { // TODO: Create delay task
//				if err = billFacade.createTransfer(ctx, splitID, memberTotal.BillIDs, memberID, currency, memberTotal.Balance()); err != nil {
//					return err
//				}
//			}
//		}
//	}
//	return nil
//}
//
//func (billFacade) createTransfer(_ context.Context, splitID int64, billIDs []int64, memberID, currency string, amount decimal.Decimal64p2) error {
//	return nil
//}
