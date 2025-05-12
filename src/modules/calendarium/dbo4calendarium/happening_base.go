package dbo4calendarium

import (
	"fmt"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

type HappeningBase struct {
	Type         HappeningType             `json:"type" firestore:"type"`
	Status       string                    `json:"status" firestore:"status"`
	Cancellation *Cancellation             `json:"canceled,omitempty" firestore:"canceled,omitempty"`
	Title        string                    `json:"title" firestore:"title"`
	Summary      string                    `json:"summary,omitempty" firestore:"summary,omitempty"`
	Levels       []string                  `json:"levels,omitempty" firestore:"levels,omitempty"`
	Slots        map[string]*HappeningSlot `json:"slots,omitempty" firestore:"slots,omitempty"`
	WithHappeningPrices
}

func (v *HappeningBase) MarkAsCanceled(cancellation Cancellation) (updates []update.Update) {
	if v.Status != HappeningStatusCanceled {
		v.Status = HappeningStatusCanceled
		updates = append(updates, update.ByFieldName("status", v.Status))
	}
	if v.Cancellation == nil {
		v.Cancellation = &cancellation
		updates = append(updates, update.ByFieldName("canceled", v.Cancellation))
	}
	return
}

func (v *HappeningBase) RemoveCancellation() (updates []update.Update) {
	if v.Status == HappeningStatusCanceled {
		v.Status = HappeningStatusActive
		updates = append(updates, update.ByFieldName("status", v.Status))
	}
	if v.Cancellation != nil {
		v.Cancellation = nil
		updates = append(updates, update.ByFieldName("canceled", update.DeleteField))
	}
	return updates
}

func (v *HappeningBase) GetSlot(id string) (slot *HappeningSlot) {
	return v.Slots[id]
}

func (v *HappeningBase) HasSlot(id string) bool {
	_, ok := v.Slots[id]
	return ok
}

// Validate returns error if not valid
func (v *HappeningBase) Validate() error {
	switch v.Type {
	case HappeningTypeSingle, HappeningTypeRecurring:
		break
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	default:
		return validation.NewErrBadRecordFieldValue("type", "unknown value: "+v.Type)
	}
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	if len(v.Title) > 100 {
		return validation.NewErrBadRequestFieldValue("title", "too long, max 100 characters")
	}
	if len(v.Summary) > 200 {
		return validation.NewErrBadRequestFieldValue("summary", "too long, max 200 characters")
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
	if v.Cancellation != nil && v.Status != HappeningStatusCanceled && v.Status != HappeningStatusDeleted {
		return validation.NewErrBadRecordFieldValue("canceled", "can be populated only for canceled or deleted happenings, current status="+v.Status)
	}

	if err := dbmodels.ValidateTitle(v.Title); err != nil {
		return err
	}
	if len(v.Slots) == 0 {
		return validation.NewErrRecordIsMissingRequiredField(SlotsField)
	}

	for slotID, slot := range v.Slots {
		if strings.TrimSpace(slotID) == "" {
			return validation.NewErrBadRecordFieldValue(SlotsField, "has empty string key")
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
	return nil
}
