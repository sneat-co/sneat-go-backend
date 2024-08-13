package models4debtus

import (
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/strongo/decimal"
	"time"
)

var ErrBalanceIsZero = errors.New("balance is zero")

func updateBalanceWithInterest(failOnZeroBalance bool, b money.Balance, outstandingWithInterest []TransferWithInterestJson, periodEnds time.Time) (err error) {
	for _, outstandingTransferWithInterest := range outstandingWithInterest {
		if balanceValue := b[outstandingTransferWithInterest.Currency]; balanceValue == 0 && failOnZeroBalance {
			return fmt.Errorf("%w: outstanding transfer %v with currency %v is not presented in balance", ErrBalanceIsZero, outstandingTransferWithInterest.TransferID, outstandingTransferWithInterest.Currency)
		} else {
			interestValue := calculateInterestValue(outstandingTransferWithInterest, periodEnds)
			if balanceValue < 0 {
				interestValue = -interestValue
			}
			b[outstandingTransferWithInterest.Currency] = balanceValue + interestValue
		}
	}
	return
}

func (t *TransferData) validateTransferInterestAndReturns() (err error) {
	if err = t.TransferInterest.ValidateTransferInterest(); err != nil {
		return
	}
	if t.AmountInCentsInterest < 0 {
		panic(fmt.Sprintf("t.AmountInCentsInterest < 0: %v", t.AmountInCentsInterest))
	}
	if !t.IsReturn && t.AmountInCentsInterest != 0 {
		panic(fmt.Sprintf("!t.IsReturn && t.AmountInCentsInterest != 0: %v", t.AmountInCentsInterest))
	}
	if t.AmountInCentsInterest > t.AmountInCents {
		panic(fmt.Sprintf("t.AmountInCentsInterest > t.AmountInCents: %v > %v", t.AmountInCentsInterest, t.AmountInCents))
	}

	if t.InterestType != "" { // TODO: Migrate old records and then do the check for all api4transfers
		returns := t.GetReturns()
		var amountReturned decimal.Decimal64p2
		for _, r := range returns {
			amountReturned += r.Amount
		}
		if amountReturned != t.AmountReturned() {
			return fmt.Errorf("sum(returns.Amount) != *TransferData.AmountInCentsReturned: %v != %v", amountReturned, t.AmountInCentsReturned)
		}
	}
	return
}
