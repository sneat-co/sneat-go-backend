package dbo4calendarium

import (
	"errors"
	"fmt"
	"github.com/strongo/validation"
)

// HappeningAdjustment at the moment supposed to be used only for recurring happenings
type HappeningAdjustment struct {
	Slots map[string]*SlotAdjustment `json:"slots,omitempty" firestore:"slots,omitempty"`
}

func (v HappeningAdjustment) IsEmpty() bool {
	return len(v.Slots) == 0
}

func (v *HappeningAdjustment) Validate() error {
	if v == nil {
		return nil
	}
	for happeningID, slot := range v.Slots {
		if err := slot.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("slots[%v]", happeningID), err.Error())
		}

	}
	return nil
}

type SlotAdjustment struct {
	Adjustment   *HappeningSlot `json:"adjustment,omitempty" firestore:"adjustment,omitempty"`
	Cancellation *Cancellation  `json:"canceled,omitempty" firestore:"canceled,omitempty"`
}

func (v *SlotAdjustment) IsEmpty() bool {
	return v.Cancellation == nil && v.Adjustment.IsEmpty()
}

func (v *SlotAdjustment) Validate() error {
	if v == nil {
		return errors.New("nil")
	}

	if v.Cancellation != nil {
		if err := v.Cancellation.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("cancellation", err.Error())
		}
	}
	return nil
}
