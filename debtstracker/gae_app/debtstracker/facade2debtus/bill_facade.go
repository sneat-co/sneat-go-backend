package facade2debtus

import (
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"math"

	"strconv"
	"time"

	"context"
	"github.com/strongo/decimal"
)

type billFacade struct {
}

var Bill = billFacade{}

func (billFacade) AssignBillToGroup(c context.Context, tx dal.ReadwriteTransaction, inBill models.Bill, groupID, userID string) (bill models.Bill, group models.GroupEntry, err error) {
	bill = inBill
	if err = bill.Data.AssignToGroup(groupID); err != nil {
		return
	}
	if bill.Data.MembersCount == 0 {
		if group, err = dtdal.Group.GetGroupByID(c, tx, groupID); err != nil {
			return
		}
		if group.Data.MembersCount > 0 {
			groupMembers := group.Data.GetGroupMembers()

			billMembers := make([]models.BillMemberJson, len(groupMembers))
			paidIsSet := false
			for i, gm := range groupMembers {
				billMembers[i] = models.BillMemberJson{
					MemberJson: gm.MemberJson,
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
					//group.AddOrGetMember(userID, "", )
					var user models.AppUser
					if user, err = User.GetUserByID(c, tx, userID); err != nil {
						return
					}
					_, _, _, groupMember, _ := group.Data.AddOrGetMember(userID, "", user.Data.FullName())

					billMembers = append(billMembers, models.BillMemberJson{
						MemberJson: groupMember.MemberJson,
						Paid:       bill.Data.AmountTotal,
					})
				}
			}
			if err = bill.Data.SetBillMembers(billMembers); err != nil {
				return
			}
			if bill.Data.Currency != money.CurrencyCode("") {
				if _, err = group.Data.ApplyBillBalanceDifference(bill.Data.Currency, bill.Data.GetBalance().BillBalanceDifference(models.BillBalanceByMember{})); err != nil {
					return
				}
				if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
					return
				}
			}
			logus.Debugf(c, "bill.GetBillMembers(): %+v", bill.Data.GetBillMembers())
		}
	}
	return
}

func (billFacade) CreateBill(c context.Context, tx dal.ReadwriteTransaction, billEntity *models.BillEntity) (bill models.Bill, err error) {
	if c == nil {
		panic("Parameter c context.Context is required")
	}
	logus.Debugf(c, "billFacade.CreateBill(%v)", billEntity)
	if tx == nil {
		panic("parameter tx dal.ReadwriteTransaction is required")
	}
	if billEntity == nil {
		panic("Parameter billEntity *models.BillEntity is required")
	}
	if !models.IsValidBillSplit(billEntity.SplitMode) {
		panic(fmt.Sprintf("billEntity.SplitMode has unknown value: %v", billEntity.SplitMode))
	}
	if billEntity.CreatorUserID == "" {
		err = fmt.Errorf("%w: billEntity.CreatorUserID == 0", ErrBadInput)
		return
	}
	if billEntity.SplitMode == "" {
		err = fmt.Errorf("%w: Missing required property SplitMode", ErrBadInput)
		return
	}
	if billEntity.AmountTotal == 0 {
		err = fmt.Errorf("%w: billEntity.AmountTotal == 0", ErrBadInput)
		return
	}
	if billEntity.AmountTotal < 0 {
		err = fmt.Errorf("%w: billEntity.AmountTotal < 0: %v", ErrBadInput, billEntity.AmountTotal)
		return
	}
	if billEntity.Status == "" {
		err = fmt.Errorf("%w: billEntity.Status property is required", err)
		return
	}
	if !models.IsValidBillStatus(billEntity.Status) {
		err = fmt.Errorf("%w: invalid status: %v, expected one of %v", ErrBadInput, billEntity.Status, models.BillStatuses)
		return
	}

	billEntity.DtCreated = time.Now()

	members := billEntity.GetBillMembers()
	//if len(members) == 0 {
	//	return bill, fmt.Errorf("len(members) == 0, MembersJson: %v", billEntity.MembersJson)
	//}

	if len(members) == 0 {
		billEntity.SplitMode = models.SplitModeEqually
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
		case models.SplitModeShare:
			shareAmount = decimal.NewDecimal64p2FromFloat64(
				math.Floor(billEntity.AmountTotal.AsFloat64()/float64(len(members))*100+0.5) / 100,
			)
		case models.SplitModeEqually:
			amountToSplitEqually := billEntity.AmountTotal
			var totalAdjustmentByMembers decimal.Decimal64p2
			for i, member := range members {
				if member.Adjustment > billEntity.AmountTotal {
					err = fmt.Errorf("%w: members[%d].Adjustment > billEntity.AmountTotal", ErrBadInput, i)
					return
				} else if member.Adjustment < 0 && member.Adjustment < -1*billEntity.AmountTotal {
					err = fmt.Errorf("%w: members[%d].AdjustmentInCents < 0 && AdjustmentInCents < -1*billEntity.AmountTotal", ErrBadInput, i)
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
					err = fmt.Errorf("%w: members[%d].Owes is negative: %v", ErrBadInput, i, member.Owes)
					return
				}
				if member.UserID != billEntity.CreatorUserID {
					if len(member.ContactByUser) == 0 {
						err = fmt.Errorf("len(members[i].ContactByUser) == 0: i==%v", i)
						return
					}
					if member.UserID == "" {
						if len(member.ContactByUser) == 0 {
							err = errors.New("bill member is missing ContactByUser ID")
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
				case models.SplitModeShare:
					expectedAmount := int64(shareAmount) * int64(member.Shares)
					deviation := expectedAmount - int64(member.Owes)
					if deviation > 1 || deviation < -1 {
						return fmt.Errorf("%w: member #%d has amount %v deviated too much (for %v) from expected %v", ErrBadInput, i, member.Owes, decimal.Decimal64p2(deviation), decimal.Decimal64p2(expectedAmount))
					}
				default:
					deviation := int64(member.Owes - member.Adjustment - equalAmount)
					if deviation > 1 || deviation < -1 {
						return fmt.Errorf("%w: member #%d has amount %v deviated too much (for %v) from equal %v", ErrBadInput, i, member.Owes, decimal.Decimal64p2(deviation), equalAmount)
					}
				}
				return nil
			}
			switch billEntity.SplitMode {
			case models.SplitModeEqually:
				// ensureNoAdjustment()
				ensureEqualShare()
				if err = ensureMemberAmountDeviateWithin1cent(); err != nil {
					return
				}
				amountsCountByValue[member.Owes]++
			case models.SplitModeExactAmount:
				ensureNoAdjustment()
				ensureNoShare()
			case models.SplitModePercentage:
				totalPercentageByMembers += member.Percent
				// ensureNoAdjustment()
			case models.SplitModeShare:
				if member.Shares == 0 {
					err = fmt.Errorf("%w: member %d is missing Shares value", ErrBadInput, i)
					return
				}
				// ensureNoAdjustment()
			}
		}

		if payersCount > 1 {
			err = ErrBillHasTooManyPayers
			return
		}

		if !(billEntity.Status == models.STATUS_DRAFT && totalPaidByMembers == 0) && totalPaidByMembers != billEntity.AmountTotal {
			err = fmt.Errorf("%w: total paid for all members should be equal to billEntity amount (%d), got %d", ErrBadInput, billEntity.AmountTotal, totalPaidByMembers)
			return
		}
		switch billEntity.SplitMode {
		case models.SplitModeEqually:
			if len(amountsCountByValue) > 2+adjustmentsCount {
				err = fmt.Errorf("%w: len(amountsCountByValue):%d > 2 + adjustmentsCount:%d", ErrBadInput, amountsCountByValue, adjustmentsCount)
				return
			}
		case models.SplitModePercentage:
			if totalPercentageByMembers != decimal.NewDecimal64p2FromInt(100) {
				err = fmt.Errorf("%w: total percentage for all members should be 100%%, got %v%%", ErrBadInput, totalPercentageByMembers)
				return
			}
		case models.SplitModeShare:
			if billEntity.Shares == 0 {
				billEntity.Shares = totalSharesPerMembers
			} else if billEntity.Shares != totalSharesPerMembers {
				err = fmt.Errorf("%w: billEntity.Shares != totalSharesPerMembers", ErrBadInput)
				return
			}
		}

		if (totalOwedByMembers != 0 || totalPaidByMembers != 0) && totalOwedByMembers != billEntity.AmountTotal {
			err = fmt.Errorf("totalOwedByMembers != billEntity.AmountTotal: %v != %v", totalOwedByMembers, billEntity.AmountTotal)
			return
		}

		// Load counterparties so we can get respective userIDs
		var counterparties []models.ContactEntry
		// Use non transactional context
		counterparties, err = GetContactsByIDs(c, tx, contactIDs)
		if err != nil {
			err = fmt.Errorf("failed to get counterparties by ID: %w", err)
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

	if bill, err = InsertBillEntity(c, tx, billEntity); err != nil {
		return
	}

	billHistory := models.NewBillHistoryBillCreated(bill, nil)
	if err = dtdal.InsertWithRandomStringID(c, tx, billHistory.Record); err != nil {
		return
	}
	return
}

//func (billFacade) CreateBillTransfers(c context.Context, billID string) error {
//	bill, err := facade2debtus.GetBillByID(c, billID)
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
//			if err = Bill.createBillTransfer(c, billID, strconv.FormatInt(bill.CreatorUserID, 10)); err != nil {
//				return errors.Wrapf(err, "Failed to create bill trasfer for %d", creatorContactID)
//			}
//		}
//	}
//	return nil
//}
//
//func (billFacade) createBillTransfer(c context.Context, billID string, creatorCounterpartyID int64) error {
//	err := dtdal.DB.RunInTransaction(c, func(c context.Context) error {
//		bill, err := facade2debtus.GetBillByID(c, billID)
//
//		if err != nil {
//			return err
//		}
//		members := bill.GetBillMembers()
//
//		var (
//			borrower *models.BillMemberJson
//			payer    *models.BillMemberJson
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
//			return errors.New("Bill member not found by creatorCounterpartyID")
//		}
//		if payer == nil {
//			return ErrBillHasNoPayer
//		}
//		//transferSource := dtdal.NewTransferSourceBot("api", "no-id", "0") // TODO: Needs refactoring! Move it out of DAL, do we really need an interface?
//
//		from := models.TransferCounterpartyInfo{
//			UserID:    payer.UserID,
//			ContactID: payer.ContactByUser[sCreatorUserID].ContactID,
//		}
//		to := models.TransferCounterpartyInfo{
//			UserID:    borrower.UserID,
//			ContactID: payer.ContactByUser[sCreatorUserID].ContactID,
//		}
//		logus.Debugf(c, "from: %v", from)
//		logus.Debugf(c, "to: %v", to)
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

func (billFacade) GetBillMembersUserInfo(c context.Context, bill models.Bill, forUserID int64) (billMembersUserInfo []BillMemberUserInfo, err error) {
	sUserID := strconv.FormatInt(forUserID, 10)

	for i, member := range bill.Data.GetBillMembers() {
		var (
			billMemberContactJson models.MemberContactJson
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

func (billFacade) AddBillMember(
	c context.Context, tx dal.ReadwriteTransaction, userID string, inBill models.Bill, memberID, memberUserID string, memberUserName string, paid decimal.Decimal64p2,
) (
	bill models.Bill, group models.GroupEntry, changed, isJoined bool, err error,
) {
	if tx == nil {
		panic("This method should be called within transaction")
	}
	logus.Debugf(c, "billFacade.AddBillMember(bill.ID=%v, memberID=%v, memberUserID=%v, memberUserName=%v, paid=%v)", bill.ID, memberID, memberUserID, memberUserName, paid)
	if paid < 0 {
		panic("paid < 0")
	}
	bill = inBill
	if bill.ID == "" {
		panic("bill.ID is empty string")
	}

	// TODO: Verify bill was obtained within transaction

	previousBalance := bill.Data.GetBalance()

	var (
		//isNew bool
		index                  int
		groupChanged           bool
		groupMember            models.GroupMemberJson
		billMember             models.BillMemberJson
		billMembers            []models.BillMemberJson
		groupMembers           []models.GroupMemberJson
		groupMembersJsonBefore string
	)

	totalAboutBefore := bill.Data.AmountTotal

	if userGroupID := bill.Data.GetUserGroupID(); userGroupID != "" {
		if group, err = dtdal.Group.GetGroupByID(c, tx, userGroupID); err != nil {
			return
		}

		groupMembersJsonBefore = group.Data.MembersJson

		if _, groupChanged, _, groupMember, groupMembers = group.Data.AddOrGetMember(memberUserID, "", memberUserName); groupChanged {
			group.Data.SetGroupMembers(groupMembers)
		} else {
			logus.Debugf(c, "GroupEntry billMembers not changed, groupMember.ID: "+groupMember.ID)
		}
	}

	_, changed, index, billMember, billMembers = bill.Data.AddOrGetMember(groupMember.ID, memberUserID, "", memberUserName)

	logus.Debugf(c, "billMember.ID: "+billMember.ID)

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

	logus.Debugf(c, "billMembers: %+v", billMembers)
	if err = bill.Data.SetBillMembers(billMembers); err != nil {
		return
	}
	logus.Debugf(c, "bill.GetBillMembers(): %+v", bill.Data.GetBillMembers())

	if err = dtdal.Bill.SaveBill(c, tx, bill); err != nil {
		return
	}

	logus.Debugf(c, "bill.GetBillMembers() after save: %v", bill.Data.GetBillMembers())

	currentBalance := bill.Data.GetBalance()

	if balanceDifference := currentBalance.BillBalanceDifference(previousBalance); balanceDifference.IsNoDifference() {
		logus.Debugf(c, "Bill balanceDifference: %v", balanceDifference)
		if groupChanged, err = group.Data.ApplyBillBalanceDifference(bill.Data.Currency, balanceDifference); err != nil {
			err = fmt.Errorf("failed to apply bill difference: %w", err)
			return
		}
		if groupChanged {
			if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
				return
			}
		}
	}

	logus.Debugf(c, "group: %+v", group)
	var groupMembersJsonAfter string
	if group.Data != nil {
		groupMembersJsonAfter = group.Data.MembersJson
	}
	billHistory := models.NewBillHistoryMemberAdded(userID, bill, totalAboutBefore, groupMembersJsonBefore, groupMembersJsonAfter)

	if err = tx.Insert(c, billHistory.Record); err != nil {
		return
	}

	isJoined = true
	return
}

var (
	ErrSettledBillsCanNotBeDeleted   = errors.New("settled bills can't be deleted")
	ErrOnlyDeletedBillsCanBeRestored = errors.New("only deleted bills can be restored")
)

func (billFacade) DeleteBill(c context.Context, billID string, userID string) (bill models.Bill, err error) {
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if bill, err = GetBillByID(c, nil, billID); err != nil {
			return
		}
		if bill.Data.Status == models.BillStatusSettled {
			err = ErrSettledBillsCanNotBeDeleted
			return
		}
		if bill.Data.Status == models.BillStatusDraft || bill.Data.Status == models.BillStatusOutstanding {
			billHistory := models.NewBillHistoryBillDeleted(userID, bill)
			if err = tx.Insert(c, billHistory.Record); err != nil {
				return
			}
			bill.Data.Status = models.BillStatusDeleted
			if err = dtdal.Bill.SaveBill(c, tx, bill); err != nil {
				return
			}
		}
		if groupID := bill.Data.GetUserGroupID(); groupID != "" {
			var group models.GroupEntry
			if group, err = dtdal.Group.GetGroupByID(c, tx, groupID); err != nil {
				return
			}
			outstandingBills := group.Data.GetOutstandingBills()
			for i, billJson := range outstandingBills {
				if billJson.ID == billID {
					outstandingBills = append(outstandingBills[:i], outstandingBills[i+1:]...)
					if _, err = group.Data.SetOutstandingBills(outstandingBills); err != nil {
						return err
					}
					groupMembers := group.Data.GetGroupMembers()
					billMembers := bill.Data.GetBillMembers()
					for j, groupMember := range groupMembers {
						for _, billMember := range billMembers {
							if billMember.ID == groupMember.ID {
								groupMember.Balance[bill.Data.Currency] -= billMember.Balance()
								groupMembers[j] = groupMember
								break
							}
						}
					}
					group.Data.SetGroupMembers(groupMembers)
					if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
						return
					}
					break
				}
			}
		}
		return
	}, dal.TxWithCrossGroup()); err != nil {
		return
	}
	return
}

func (billFacade) RestoreBill(c context.Context, billID string, userID string) (bill models.Bill, err error) {
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if bill, err = GetBillByID(c, nil, billID); err != nil {
			return
		}
		if bill.Data.Status != models.BillStatusDeleted {
			err = ErrOnlyDeletedBillsCanBeRestored
			return
		}

		if bill.Data.MembersCount > 1 {
			bill.Data.Status = models.BillStatusOutstanding
		} else {
			bill.Data.Status = models.BillStatusDraft
		}
		billHistory := models.NewBillHistoryBillRestored(userID, bill)
		if err = tx.Insert(c, billHistory.Record); err != nil {
			return
		}
		if err = dtdal.Bill.SaveBill(c, tx, bill); err != nil {
			return
		}
		if groupID := bill.Data.GetUserGroupID(); groupID != "" {
			var group models.GroupEntry
			if group, err = dtdal.Group.GetGroupByID(c, tx, groupID); err != nil {
				return
			}
			var groupChanged bool
			if groupChanged, err = group.Data.AddBill(bill); err != nil {
				return
			} else if groupChanged {
				if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
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

func GetBillByID(c context.Context, tx dal.ReadSession, billID string) (bill models.Bill, err error) {
	if tx == nil {
		if tx, err = facade.GetDatabase(c); err != nil {
			return bill, err
		}
	}
	bill.ID = billID
	bill.Key = dal.NewKeyWithID(models.BillKind, billID)
	bill.Data = new(models.BillEntity)
	bill.Record = dal.NewRecordWithData(bill.Key, bill.Data)
	err = tx.Get(c, bill.Record)
	return
}

func InsertBillEntity(c context.Context, tx dal.ReadwriteTransaction, billEntity *models.BillEntity) (bill models.Bill, err error) {
	if billEntity == nil {
		panic("billEntity == nil")
	}
	if billEntity.CreatorUserID == "" {
		panic("CreatorUserID == 0")
	}
	if billEntity.AmountTotal == 0 {
		panic("AmountTotal == 0")
	}

	billEntity.DtCreated = time.Now()
	bill.Data = billEntity

	err = tx.Insert(c, bill.Record)
	return
}

//func (billFacade billFacade) createTransfers(c context.Context, splitID int64) error {
//	split, err := dtdal.Split.GetSplitByID(c, splitID)
//	if err != nil {
//		return err
//	}
//	bills, err := dtdal.Bill.GetBillsByIDs(c, split.BillIDs)
//
//	balances := billFacade.getBalances(splitID, bills)
//	balances = billFacade.cleanupBalances(balances)
//
//	for currency, totalsByMember := range balances {
//		for memberID, memberTotal := range totalsByMember {
//			if memberTotal.Balance() > 0 { // TODO: Create delay task
//				if err = billFacade.createTransfer(c, splitID, memberTotal.BillIDs, memberID, currency, memberTotal.Balance()); err != nil {
//					return err
//				}
//			}
//		}
//	}
//	return nil
//}
//
//func (billFacade) createTransfer(c context.Context, splitID int64, billIDs []int64, memberID, currency string, amount decimal.Decimal64p2) error {
//	return nil
//}
