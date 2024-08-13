package facade4splitus

import (
	"github.com/crediterra/money"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/strongo/decimal"
)

type SplitMemberTotal struct {
	briefs4splitus.BillMemberBalance
}

func (t SplitMemberTotal) Balance() decimal.Decimal64p2 {
	return t.Paid - t.Owes
}

type SplitBalanceByMember map[string]SplitMemberTotal
type SplitBalanceByCurrencyAndMember map[money.CurrencyCode]SplitBalanceByMember

//func (billFacade) getBalances(splitID int64, bills []models.Bill) (balanceByCurrency SplitBalanceByCurrencyAndMember) {
//	balanceByCurrency = make(SplitBalanceByCurrencyAndMember)
//	for _, bill := range bills {
//		var (
//			balanceByMember SplitBalanceByMember
//			ok              bool
//		)
//		if balanceByMember, ok = balanceByCurrency[bill.Data.Currency]; !ok {
//			balanceByMember = make(SplitBalanceByMember)
//			balanceByCurrency[bill.Data.Currency] = balanceByMember
//		}
//		for memberID, memberBalance := range bill.Data.GetBalance() {
//			memberTotal := balanceByMember[memberID]
//			memberTotal.Paid += memberBalance.Paid
//			memberTotal.Owes += memberBalance.Owes
//			balanceByMember[memberID] = memberTotal
//		}
//	}
//	return
//}

//func (billFacade) cleanupBalances(balanceByCurrency SplitBalanceByCurrencyAndMember) SplitBalanceByCurrencyAndMember {
//	return balanceByCurrency
//}
