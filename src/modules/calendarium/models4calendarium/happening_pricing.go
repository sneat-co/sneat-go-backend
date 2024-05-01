package models4calendarium

import (
	"github.com/crediterra/money"
	"github.com/strongo/validation"
	"strconv"
)

// HappeningPrices describes prices for happening
type HappeningPrices struct {
	Prices []HappeningPrice `json:"prices,omitempty" firestore:"prices,omitempty"`
}

// HappeningPrice describes price for happening
type HappeningPrice struct {
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
