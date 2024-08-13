package maintainance

import (
	"github.com/crediterra/money"
	"strings"
)

func FixBalanceCurrencies(balance money.Balance) (changed bool) {
	euro := money.CurrencyCode("euro")
	for c, v := range balance {
		if c == euro {
			c = money.CurrencyEUR
		}
		if len(c) == 3 {
			cc := strings.ToUpper(string(c))
			if cc != string(c) {
				if cu := money.CurrencyCode(cc); cu.IsMoney() {
					balance[cu] += v
					delete(balance, c)
					changed = true
				}
			}
		}
	}
	return
}
