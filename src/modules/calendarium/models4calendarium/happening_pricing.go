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

// GetPriceByID returns price by ID
func (v WithHappeningPrices) GetPriceByID(priceID string) *HappeningPrice {
	for _, price := range v.Prices {
		if price.ID == priceID {
			return price
		}
	}
	return nil
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
		id := price.Term.ID()
		for j, p := range v.Prices {
			if i != j && p.ID == id {
				return validation.NewErrBadRecordFieldValue("prices",
					fmt.Sprintf("duplicate price ID at indexes %d & %d: %s", i, j, id))
			}
		}
	}
	return nil
}

// HappeningPrice describes price for happening
type HappeningPrice struct {
	ID     string       `json:"id,omitempty" firestore:"id,omitempty"`
	Term   Term         `json:"term" firestore:"term"`
	Amount money.Amount `json:"amount" firestore:"amount"`

	// Does not need to be *int as '0' means not applicable and is omitted from JSON & Firestore
	ExpenseQuantity int `json:"expenseQuantity,omitempty" firestore:"expenseQuantity,omitempty"`
}

// Validate returns error if not valid
func (v HappeningPrice) Validate() error {
	// We do not validate ID for "" here as it will be empty in request to create a new price entry
	// Though we validate it's not empty in HappeningBrief
	if v.ID == "*" {
		return validation.NewErrBadRecordFieldValue("id", "should not be '*'")
	}
	if err := v.Term.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("term", err.Error())
	}
	if err := v.Amount.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("amount", err.Error())
	}
	if v.Amount.Value < 0 {
		return validation.NewErrBadRecordFieldValue("amount", "should be positive, got: "+v.Amount.String())
	}
	if v.ExpenseQuantity < 0 {
		return validation.NewErrBadRecordFieldValue("expenseQuantity", "should be positive or zero, got: "+strconv.Itoa(v.ExpenseQuantity))
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

// ID returns unique identifier of the term
func (v Term) ID() string {
	return fmt.Sprintf("%s%d", v.Unit, v.Length)
}

// String returns string representation of the term
func (v Term) String() string {
	if v.Unit == TermUnitSingle {
		return "single"
	}
	if v.Length == 1 {
		return fmt.Sprintf("1 %s", v.Unit)
	}
	return fmt.Sprintf("%d %ss", v.Length, v.Unit)
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
