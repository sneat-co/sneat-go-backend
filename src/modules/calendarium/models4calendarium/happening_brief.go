package models4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// HappeningBrief hold data that stored both in entity record and in a brief.
type HappeningBrief struct {
	Type     HappeningType    `json:"type" firestore:"type"`
	Status   string           `json:"status" firestore:"status"`
	Canceled *Canceled        `json:"canceled,omitempty" firestore:"canceled,omitempty"`
	Title    string           `json:"title" firestore:"title"`
	Levels   []string         `json:"levels,omitempty" firestore:"levels,omitempty"`
	Slots    []*HappeningSlot `json:"slots,omitempty" firestore:"slots,omitempty"`
	WithHappeningPrices
}

func (v HappeningBrief) GetSlot(id string) (i int, slot *HappeningSlot) {
	for i, slot = range v.Slots {
		if slot.ID == id {
			return
		}
	}
	return -1, nil
}

// Validate returns error if not valid
func (v HappeningBrief) Validate() error {
	switch v.Type {
	case HappeningTypeSingle, HappeningTypeRecurring:
		break
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	default:
		return validation.NewErrBadRecordFieldValue("type", "unknown value: "+v.Type)
	}
	if v.Status == "" {
		return validation.NewErrRecordIsMissingRequiredField("status")
	}
	if !IsKnownHappeningStatus(v.Status) {
		return validation.NewErrBadRecordFieldValue("status", fmt.Sprintf("unknown value: '%v'", v.Status))
	}
	if v.Status == HappeningStatusCanceled && v.Canceled == nil {
		return validation.NewErrRecordIsMissingRequiredField("canceled")
	}
	if v.Canceled != nil && v.Status != HappeningStatusCanceled {
		return validation.NewErrBadRecordFieldValue("canceled", "should be populated only for canceled happenings, current status="+v.Status)
	}

	if err := dbmodels.ValidateTitle(v.Title); err != nil {
		return err
	}
	if len(v.Slots) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("slots")
	}

	for i, slot := range v.Slots {
		if err := slot.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("slots[%v]", i), err.Error())
		}
		for j, s := range v.Slots {
			if i != j && s.ID == slot.ID {
				return validation.NewErrBadRecordFieldValue("slots", fmt.Sprintf("at least 2 slots have same ContactID at indexes: %v & %v", i, j))
			}
			// TODO: Add more validations?
		}
	}
	if err := v.WithHappeningPrices.Validate(); err != nil {
		return err // No need for field name as it reported by WithHappeningPrices.Validate()
	}
	for i, price := range v.Prices {
		if price.ID == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("prices[%v].id", i), "empty string value")
		}
	}
	return nil
}
