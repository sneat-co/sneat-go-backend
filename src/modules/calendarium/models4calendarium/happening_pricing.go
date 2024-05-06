package models4calendarium

import (
	"fmt"
	"github.com/crediterra/money"
	"github.com/strongo/validation"
	"strconv"
)

// WithHappeningPrices describes prices for happening
type WithHappeningPrices struct {
	Prices []*HappeningPrice `json:"prices,omitempty" firestore:"prices,omitempty"`
}

// Validate returns error if not valid
func (v WithHappeningPrices) Validate() error {
	for i, price := range v.Prices {
		if price == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("prices[%d]", i), "nil value")
		}
		if err := price.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("prices["+strconv.Itoa(i)+"]", err.Error())
		}
	}
	return nil
}

// HappeningPrice describes price for happening
type HappeningPrice struct {
	ID     string       `json:"id,omitempty" firestore:"id,omitempty"`
	Term   Term         `json:"term" firestore:"term"`
	Amount money.Amount `json:"amount" firestore:"amount"`
}

// Validate returns error if not valid
func (v HappeningPrice) Validate() error {
	if err := v.Term.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("term", err.Error())
	}
	if err := v.Amount.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("amount", err.Error())
	}
	if v.Amount.Value < 0 {
		return validation.NewErrBadRecordFieldValue("amount", "should be positive, got: "+v.Amount.String())
	}
	return nil
}

// TermUnit describes unit of the term
type TermUnit string

const (
	TermUnitSingle TermUnit = "single"
	TermUnitSecond TermUnit = "second"
	TermUnitMinute TermUnit = "minute"
	TermUnitHour   TermUnit = "hour"
	TermUnitDay    TermUnit = "day"
	TermUnitWeek   TermUnit = "week"
	TermUnitMonth  TermUnit = "month"
	TermUnitYear   TermUnit = "year"
)

// Term describes duration of the price
type Term struct {
	Unit   TermUnit `json:"unit" firestore:"unit"`
	Length int      `json:"length" firestore:"length"`
}

// Validate returns error if not valid
func (v Term) Validate() error {
	switch v.Unit {
	case
		TermUnitSingle,
		TermUnitSecond,
		TermUnitMinute,
		TermUnitHour,
		TermUnitDay,
		TermUnitWeek,
		TermUnitMonth,
		TermUnitYear:
		break
	case "":
		return validation.NewErrRecordIsMissingRequiredField("term")
	default:
		return validation.NewErrBadRecordFieldValue("term", "unknown value: "+string(v.Unit))
	}
	if v.Length == 0 {
		return validation.NewErrRecordIsMissingRequiredField("length")
	} else if v.Length < 0 {
		return validation.NewErrBadRecordFieldValue("length", "should be positive, got: "+strconv.Itoa(v.Length))
	}
	return nil
}
