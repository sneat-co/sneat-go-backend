package models

import (
	"fmt"

	"github.com/strongo/decimal"
	"time"

	"github.com/crediterra/go-interest"
	"github.com/crediterra/money"
	"github.com/sanity-io/litter"
)

type TransferInterest struct {
	InterestType          interest.Formula          `datastore:",noindex,omitempty"`
	InterestPeriod        interest.RatePeriodInDays `datastore:",noindex,omitempty"`
	InterestPercent       decimal.Decimal64p2       `datastore:",noindex,omitempty"`
	InterestGracePeriod   int                       `datastore:",noindex,omitempty" json:",omitempty"` // How many days are without any interest
	InterestMinimumPeriod int                       `datastore:",noindex,omitempty" json:",omitempty"` // Minimum days for interest (e.g. penalty for earlier return).
	// InterestAmountInCents decimal.Decimal64p2 `datastore:",noindex" json:",omitempty"`
}

func NewInterest(formula interest.Formula, percent decimal.Decimal64p2, period interest.RatePeriodInDays) TransferInterest {
	if percent <= 0 {
		panic(fmt.Sprintf("percent <= 0 (%v)", percent))
	}
	if period < 0 {
		panic(fmt.Sprintf("period < 0 (%v)", period))
	}
	if !interest.IsKnownFormula(formula) {
		panic(fmt.Sprintf("unknown interest percent type (%v)", formula))
	}
	return TransferInterest{
		InterestPeriod:  period,
		InterestType:    formula,
		InterestPercent: percent,
	}
}

func (ti TransferInterest) WithMinimumPeriod(minimumPeriod int) TransferInterest {
	ti.InterestMinimumPeriod = minimumPeriod
	return ti
}

func (ti TransferInterest) WithGracePeriod(gracePeriodInDays int) TransferInterest {
	ti.InterestGracePeriod = gracePeriodInDays
	return ti
}

func NoInterest() TransferInterest {
	return TransferInterest{}
}

func (ti TransferInterest) HasInterest() bool {
	return ti.InterestPercent != 0
}

func (ti TransferInterest) GetInterestData() TransferInterest {
	return ti
}

func (ti TransferInterest) ValidateTransferInterest() (err error) {
	if ti.InterestPeriod == 0 && ti.InterestPercent == 0 && ti.InterestType == "" {
		return
	}
	if ti.InterestPeriod < 0 {
		return fmt.Errorf("InterestPeriod < 0: %v", ti.InterestPeriod)
	}
	if ti.InterestPercent <= 0 {
		return fmt.Errorf("InterestPercent <= 0: %v", ti.InterestPercent)
	}
	// if entity.InterestAmountInCents < 0 {
	// 	return fmt.Errorf("InterestAmountInCents < 0: %v", entity.InterestAmountInCents)
	// }
	if ti.InterestType == "" {
		return interest.ErrFormulaIsNotSet
	}
	if !interest.IsKnownFormula(ti.InterestType) {
		return fmt.Errorf("unknown InterestType: %v", ti.InterestType)
	}
	if ti.InterestPeriod == 0 || ti.InterestPercent == 0 {
		return fmt.Errorf(
			"one of values is 0: InterestPeriod=%v, InterestPercent=%v",
			ti.InterestPeriod,
			ti.InterestPercent,
		)
	}
	return
}

// func init() {
// 	addInterestPropertiesToClean := func(props2clean map[string]gaedb.IsOkToRemove) {
// 		props2clean["InterestType"] = gaedb.IsEmptyString
// 		props2clean["InterestPeriod"] = gaedb.IsZeroInt
// 		props2clean["InterestPercent"] = gaedb.IsZeroInt
// 		props2clean["InterestGracePeriod"] = gaedb.IsZeroInt
// 		props2clean["InterestMinimumPeriod"] = gaedb.IsZeroInt
// 	}
// 	addInterestPropertiesToClean(transferPropertiesToClean)
// }

func (t *TransferData) GetOutstandingValue(periodEnds time.Time) (outstandingValue decimal.Decimal64p2) {
	if t.IsReturn && (t.AmountInCentsReturned == 0 || t.AmountInCents == t.AmountInCentsReturned) {
		/*
			TODO: What if transfer was a return > then outstanding value? We decided to allow it for returns without interest.
		*/
		return 0
	}
	interestValue := t.GetInterestValue(periodEnds)
	outstandingValue = t.AmountInCents + interestValue - t.AmountInCentsReturned

	if outstandingValue < 0 /* && interestValue != 0 */ {
		panic(fmt.Sprintf("outstandingValue < 0: %v, IsReturn: %v, Amount: %v, Returned: %v, Interest: %v\n%v",
			outstandingValue, t.IsReturn, t.AmountInCents, t.AmountInCentsReturned, interestValue, litter.Sdump(t)))
	}
	return
}

func (t *TransferData) GetOutstandingAmount(periodEnds time.Time) money.Amount {
	return money.Amount{Currency: t.Currency, Value: t.GetOutstandingValue(periodEnds)}
}

type TransferInterestCalculable interface {
	GetLendingValue() decimal.Decimal64p2
	GetStartDate() time.Time
	GetReturns() TransferReturns
	GetInterestData() TransferInterest
}

func (t *TransferData) GetInterestValue(periodEnds time.Time) (interestValue decimal.Decimal64p2) {
	if t.InterestType == "" && t.InterestPeriod == 0 {
		return 0
	}
	return calculateInterestValue(t, periodEnds)
}

func interestDeal(t TransferInterestCalculable) interest.Deal {
	ti := t.GetInterestData()
	return interest.NewDeal(ti.InterestType, t.GetStartDate(), t.GetLendingValue(), ti.InterestPercent, interest.RatePeriodInDays(ti.InterestPeriod), ti.InterestMinimumPeriod, ti.InterestGracePeriod)
}

func interestPayments(t TransferInterestCalculable) []interest.Payment {
	returns := t.GetReturns()
	payments := make([]interest.Payment, len(returns))
	for i, r := range returns {
		payments[i] = interest.NewPayment(r.Time, r.Amount)
	}
	return payments
}

func calculateInterestValue(t TransferInterestCalculable, reportTime time.Time) (interestValue decimal.Decimal64p2) {
	deal := interestDeal(t)
	var err error
	payments := interestPayments(t)
	interestValue, _, err = interest.Calculate(reportTime, deal, payments)
	if err != nil {
		panic(err)
	}
	return
}

func (t *TransferData) AgeInDays() int {
	return interest.AgeInDays(time.Now(), t.DtCreated)
}

/*
Example:

7% per week min 3 days
1.5% в неделю мин 3 дня

*/
