package dbo4logist

import (
	"fmt"
	"strings"

	"github.com/strongo/validation"
)

type ShippingPointTask = string

const (
	ShippingPointTaskPick   ShippingPointTask = "pick"   // Pickup container
	ShippingPointTaskDrop   ShippingPointTask = "drop"   // Leave container with all the goods
	ShippingPointTaskLoad   ShippingPointTask = "load"   // Load container with goods
	ShippingPointTaskUnload ShippingPointTask = "unload" // Unload goods from container
)

type validating int // TODO move to github.com/strongo/validation package and expose as public API

const (
	ValidatingRequest validating = iota
	ValidatingRecord
)

func ValidateShippingPointTask(v ShippingPointTask, raise validating, field func() string) error {
	if strings.TrimSpace(v) != v {
		return validation.NewErrBadRecordFieldValue(field(), "must not contain leading or trailing spaces")
	}
	switch v {
	case // Valid values
		ShippingPointTaskPick,
		ShippingPointTaskDrop,
		ShippingPointTaskLoad,
		ShippingPointTaskUnload:
		return nil
	case "":
		if raise == ValidatingRequest {
			return validation.NewErrRequestIsMissingRequiredField(field())
		}
		return validation.NewErrRecordIsMissingRequiredField(field())
	default:

		m := fmt.Sprintf("unknown value: [%v]", v)
		if raise == ValidatingRequest {
			return validation.NewErrBadRequestFieldValue(field(), m)
		}
		return validation.NewErrBadRecordFieldValue(field(), m)
	}

}
func validateShippingPointTasks(v []ShippingPointTask, raise validating) error {
	for i, s := range v {
		if err := ValidateShippingPointTask(s, raise, func() string { return fmt.Sprintf("tasks[%v]", i) }); err != nil {
			return err
		}
	}
	return nil
}

// ValidateShippingPointTasksRequest validates shipping point tasks
func ValidateShippingPointTasksRequest(v []ShippingPointTask, mustHave bool) error {
	if mustHave && len(v) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("tasks")
	}
	return validateShippingPointTasks(v, ValidatingRequest)
}

// ValidateShippingPointTasksRecord validates shipping point tasks
func ValidateShippingPointTasksRecord(v []ShippingPointTask) error {
	if len(v) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("tasks")
	}
	return validateShippingPointTasks(v, ValidatingRecord)
}
