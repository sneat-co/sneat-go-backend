package dto4calendarium

import (
	"github.com/strongo/validation"
	"strings"
)

// HappeningSlotRefRequest refers to a happening slot
type HappeningSlotRefRequest struct {
	HappeningRequest
	SlotID  string `json:"slotID"`
	Weekday string `json:"weekday,omitempty"`
}

func (v HappeningSlotRefRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.SlotID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("slotID")
	}
	return nil
}
