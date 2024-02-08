package models

//go:generate ffjson $GOFILE

import (
	"github.com/crediterra/money"
	"github.com/strongo/decimal"
)

type BillJson struct {
	ID           string
	GroupID      string              `json:"g,omitempty"`
	Name         string              `json:"n"`
	MembersCount int                 `json:"m"`
	Total        decimal.Decimal64p2 `json:"t"`
	Currency     money.CurrencyCode  `json:"c"`
	UserBalance  decimal.Decimal64p2 `json:"u,omitempty"`
}

type BillMemberBalance struct {
	Paid decimal.Decimal64p2
	Owes decimal.Decimal64p2
}

func (t BillMemberBalance) Balance() decimal.Decimal64p2 {
	return t.Paid - t.Owes
}

type BillBalanceByMember map[string]BillMemberBalance

type BalanceByMember map[string]decimal.Decimal64p2

type BillBalanceDifference BalanceByMember

func (diff BillBalanceDifference) Reverse() BillBalanceDifference {
	reversed := make(BillBalanceDifference, len(diff))
	for k, v := range diff {
		reversed[k] = -v
	}
	return reversed
}

func (diff BillBalanceDifference) IsNoDifference() bool {
	diff.clear()
	return len(diff) > 0
}

func (diff BillBalanceDifference) clear() (changed bool) {
	for k, v := range diff {
		if v == 0 {
			delete(diff, k)
			changed = true
		}
	}
	return
}

func (current BillBalanceByMember) BillBalanceDifference(previous BillBalanceByMember) (difference BillBalanceDifference) {
	capacity := len(current) + 1
	if len(previous) > capacity {
		capacity = len(previous) + 1
	}

	difference = make(BillBalanceDifference, capacity)

	for memberID, mCurrent := range current {
		if diff := mCurrent.Balance() - previous[memberID].Balance(); diff != 0 {
			difference[memberID] = diff
		}
	}

	for memberID, mPrevious := range previous {
		if _, ok := current[memberID]; !ok {
			difference[memberID] = -mPrevious.Balance()
		}
	}

	return
}

type BillSettlementJson struct {
	BillID    string              `json:"bill"`
	GroupID   string              `json:"group"`
	DebtorID  string              `json:"debtor,omitempty"`
	SponsorID string              `json:"sponsor,omitempty"`
	Amount    decimal.Decimal64p2 `json:"amount,omitempty"`
}
