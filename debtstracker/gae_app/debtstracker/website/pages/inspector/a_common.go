package inspector

import (
	"github.com/crediterra/money"
	"github.com/strongo/decimal"
	"sync"
)

type balanceRow struct {
	// TODO: rename
	user                  decimal.Decimal64p2
	contacts              decimal.Decimal64p2
	transfers             decimal.Decimal64p2
	userContactBalanceErr error
	contactBalanceErr     error
}

type balancesByCurrency struct {
	*sync.Mutex
	err        error
	byCurrency map[money.CurrencyCode]balanceRow
}

type balances struct {
	withInterest    balancesByCurrency
	withoutInterest balancesByCurrency
}

func newBalances(who string, withoutInterest, withInterest money.Balance) balances {
	return balances{
		withoutInterest: newBalanceSummary(who, withoutInterest),
		withInterest:    newBalanceSummary(who, withInterest),
	}
}
