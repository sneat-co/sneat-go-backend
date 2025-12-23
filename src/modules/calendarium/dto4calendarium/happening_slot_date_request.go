package dto4calendarium

import (
	"strings"

	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

// HappeningSlotDateRequest updates slot
type HappeningSlotDateRequest struct {
	HappeningRequest
	Date string              `json:"date"`
	Slot HappeningSlotWithID `json:"slot"`
}

// Validate returns error if not valid
func (v HappeningSlotDateRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.Date) == "" {
		return validation.NewErrRecordIsMissingRequiredField("date")
	}
	if _, err := validate.DateString(v.Date); err != nil {
		return validation.NewErrBadRequestFieldValue("date", err.Error())
	}
	if err := v.Slot.Validate(); err != nil {
		return err
	}
	return nil
}

type HappeningDateSlotIDRequest struct {
	HappeningRequest
	Date   string `json:"date"`
	SlotID string `json:"slotID"`
}
