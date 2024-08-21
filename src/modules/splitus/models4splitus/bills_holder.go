package models4splitus

import (
	"errors"
	"github.com/crediterra/money"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"strings"
)

//func init() {
//	userPropertiesToClean["OutstandingBillsJson"] = gaedb.IsEmptyJSON
//	groupPropertiesToClean["OutstandingBillsJson"] = gaedb.IsEmptyJSON
//}

type BillsHolder struct { // TODO: Move out of auth package
	OutstandingBills map[string]briefs4splitus.BillBrief `json:"outstandingBills,omitempty" firestore:"outstandingBills,omitempty"`

	// Deprecated: use OutstandingBills instead
	OutstandingBillsCount int `firestore:",omitempty"`

	// Deprecated: use OutstandingBills instead
	OutstandingBillsJson string `firestore:",omitempty"`
}

func (v *BillsHolder) GetOutstandingBills() (outstandingBills map[string]briefs4splitus.BillBrief) {
	return v.OutstandingBills
}

func (v *BillsHolder) SetOutstandingBills(outstandingBills map[string]briefs4splitus.BillBrief) (err error) {
	v.OutstandingBills = outstandingBills
	return
}

func (v *BillsHolder) AddBill(bill BillEntry) (changed bool, err error) {
	outstandingBills := v.GetOutstandingBills()

	if billBrief, ok := outstandingBills[bill.ID]; ok {
		if billBrief.Name != bill.Data.Name {
			billBrief.Name = bill.Data.Name
			changed = true
		}
		if membersCount := len(bill.Data.Members); billBrief.MembersCount != membersCount {
			billBrief.MembersCount = membersCount
			changed = true
		}
		if billBrief.Total != bill.Data.AmountTotal {
			billBrief.Total = bill.Data.AmountTotal
			changed = true
		}
		goto addedOrUpdatedOrNotChanged
	}
	outstandingBills[bill.ID] = briefs4splitus.BillBrief{
		Name:         bill.Data.Name,
		MembersCount: len(bill.Data.Members),
		Total:        bill.Data.AmountTotal,
		Currency:     bill.Data.Currency,
	}
addedOrUpdatedOrNotChanged:
	if changed {
		if err = v.SetOutstandingBills(outstandingBills); err != nil {
			return
		}
		panic("not implemented yet")
		//groupMembers := v.GetGroupMembers()
		//billMembers := bill.Data.GetBillMembers()
		//for j, groupMember := range groupMembers {
		//	for _, billMember := range billMembers {
		//		if billMember.ContactID == groupMember.ContactID {
		//			groupMember.Balance[bill.Data.Currency] += billMember.Balance()
		//			groupMembers[j] = groupMember
		//			break
		//		}
		//	}
		//}
		//v.SetGroupMembers(groupMembers)
	}
	return
}

func (v *BillsHolder) ApplyBillBalanceDifference(currency money.CurrencyCode, diff briefs4splitus.BillBalanceDifference) (changed bool, err error) {
	if currency == "" {
		panic("currency parameter is required")
	}
	if strings.TrimSpace(string(currency)) != string(currency) {
		panic("currency parameter has leading ot closing spaces: " + currency)
	}

	return false, errors.New("ApplyBillBalanceDifference is not implemented yet")
	//groupMembers := v.GetGroupMembers()
	//
	//var diffTotal, balanceTotal decimal.Decimal64p2
	//diffCopy := make(models4debtus.BillBalanceDifference, len(diff))
	//
	//for i := range groupMembers {
	//	groupMemberID := groupMembers[i].ContactID
	//
	//	if memberDifference, ok := diff[groupMemberID]; ok {
	//		delete(diff, groupMemberID)
	//		diffCopy[groupMemberID] = memberDifference
	//		if memberDifference == 0 {
	//			panic("memberDifference.Paid == 0 && memberDifference.Owes == 0, memberID: " + groupMemberID)
	//		}
	//		diffTotal += memberDifference
	//		if diffAmount := memberDifference; diffAmount != 0 {
	//			if groupMembers[i].Balance == nil || len(groupMembers) == 0 {
	//				groupMembers[i].Balance = money.Balance{currency: diffAmount}
	//				balanceTotal += diffAmount
	//			} else {
	//				groupMembers[i].Balance[currency] += diffAmount
	//				if len(groupMembers[i].Balance) == 0 {
	//					groupMembers[i].Balance = nil
	//				} else {
	//					balanceTotal += groupMembers[i].Balance[currency]
	//				}
	//			}
	//		}
	//	}
	//}
	//
	//if len(diff) > 0 {
	//	err = fmt.Errorf("%w: %v", models4debtus.ErrNonGroupMember, diff)
	//	return
	//}
	//
	//if diffTotal != 0 {
	//	err = fmt.Errorf("%w: diffTotal=%v, diff=%v", models4debtus.ErrBillOwesDiffTotalIsNotZero, diffTotal, diffCopy)
	//	return
	//}
	//
	//if balanceTotal != 0 {
	//	err = fmt.Errorf("%wbalanceTotal=%v, diff=%v", models4debtus.ErrGroupTotalBalanceHasNonZeroValue, balanceTotal, diffCopy)
	//	return
	//}
	//return v.SetGroupMembers(groupMembers), err
}

func (v *BillsHolder) GetOutstandingBalance() (balance money.Balance) {
	balance = make(money.Balance, 2)
	for _, bill := range v.GetOutstandingBills() {
		balance[bill.Currency] += bill.UserBalance
	}
	return
}
