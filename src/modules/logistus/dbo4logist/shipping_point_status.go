package dbo4logist

import (
	"fmt"

	"github.com/strongo/validation"
)

type ShippingPointStatus = string

const (
	ShippingPointStatusPending    ShippingPointStatus = "pending"
	ShippingPointStatusProcessing ShippingPointStatus = "processing"
	ShippingPointStatusCompleted  ShippingPointStatus = "completed"
)

func validateShippingPointStatus(field string, v ShippingPointStatus) error {
	switch v {
	case ShippingPointStatusPending, ShippingPointStatusProcessing, ShippingPointStatusCompleted:
		return nil // OK
	case "":
		return validation.NewErrRecordIsMissingRequiredField(field)
	default:
		return validation.NewErrBadRecordFieldValue(field, fmt.Sprintf("unsupported value: [%v]", v))
	}
}
