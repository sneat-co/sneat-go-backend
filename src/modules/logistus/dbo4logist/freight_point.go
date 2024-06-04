package dbo4logist

import (
	"github.com/strongo/slice"
	"github.com/strongo/validation"
)

type FreightPoint struct {
	Tasks    []ShippingPointTask `json:"tasks" firestore:"tasks,omitempty"`
	ToLoad   *FreightLoad        `json:"toLoad,omitempty" firestore:"toLoad,omitempty"`
	ToUnload *FreightLoad        `json:"toUnload,omitempty" firestore:"toUnload,omitempty"`
}

func (v FreightPoint) HasTask(task ShippingPointTask) bool {
	return slice.Index(v.Tasks, task) >= 0
}

// Validate returns an error if the ShippingPointBase is invalid.
func (v FreightPoint) Validate() error {
	if v.ToLoad != nil {
		if err := v.ToLoad.Validate(); err != nil {
			return err // Do not wrap error here
		}
	}
	if v.ToUnload != nil {
		if err := v.ToUnload.Validate(); err != nil {
			return err // Do not wrap error here
		}
	}
	if err := ValidateShippingPointTasksRecord(v.Tasks); err != nil {
		return err
	}
	if v.ToLoad != nil && slice.Index(v.Tasks, "load") < 0 {
		return validation.NewErrBadRecordFieldValue("tasks", "must contain 'load' when toLoad is set")
	}
	if v.ToUnload != nil && slice.Index(v.Tasks, "unload") < 0 {
		return validation.NewErrBadRecordFieldValue("tasks", "must contain 'unload' when toUnload is set")
	}
	return nil
}
