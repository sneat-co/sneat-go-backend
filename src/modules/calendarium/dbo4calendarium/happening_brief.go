package dbo4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

// HappeningBrief hold data that stored both in entity record and in a brief.
type HappeningBrief struct {
	Type         HappeningType             `json:"type" firestore:"type"`
	Status       string                    `json:"status" firestore:"status"`
	Cancellation *Cancellation             `json:"canceled,omitempty" firestore:"canceled,omitempty"`
	Title        string                    `json:"title" firestore:"title"`
	Levels       []string                  `json:"levels,omitempty" firestore:"levels,omitempty"`
	Slots        map[string]*HappeningSlot `json:"slots,omitempty" firestore:"slots,omitempty"`
	WithHappeningPrices
	dbo4linkage.WithRelated
}

func (v HappeningBrief) GetSlot(id string) (slot *HappeningSlot) {
	return v.Slots[id]
}

func (v HappeningBrief) HasSlot(id string) bool {
	_, ok := v.Slots[id]
	return ok
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
	if v.Status == HappeningStatusCanceled && v.Cancellation == nil {
		return validation.NewErrRecordIsMissingRequiredField("canceled")
	}
	if v.Cancellation != nil && v.Status != HappeningStatusCanceled {
		return validation.NewErrBadRecordFieldValue("canceled", "should be populated only for canceled happenings, current status="+v.Status)
	}

	if err := dbmodels.ValidateTitle(v.Title); err != nil {
		return err
	}
	if len(v.Slots) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("slots")
	}

	for slotID, slot := range v.Slots {
		if strings.TrimSpace(slotID) == "" {
			return validation.NewErrBadRecordFieldValue("slots", "has empty string key")
		}
		if err := slot.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("slots[%v]", slotID), err.Error())
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
	if err := v.WithRelated.Validate(); err != nil {
		return err
	}
	return nil
}
