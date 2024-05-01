package models4budgetus

import (
	"github.com/crediterra/money"
	"github.com/strongo/validation"
)

type ObligationDirection string

const (
	ObligationDirectionCollection    ObligationDirection = "collection"
	ObligationDirectionDisbursements ObligationDirection = "disbursement"
)

type Obligation struct {
	Direction ObligationDirection
	Amount    money.Amount
}

func (v *Obligation) Validate() error {
	switch v.Direction {
	case ObligationDirectionCollection, ObligationDirectionDisbursements:
		break
	case "":
		return validation.NewErrRecordIsMissingRequiredField("direction")
	default:
		return validation.NewErrBadRecordFieldValue("direction", "unknown value: "+string(v.Direction))
	}
	if err := v.Amount.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("amount", err.Error())
	}
	return nil
}
