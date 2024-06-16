package dto4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/strongo/validation"
	"strings"
)

// HappeningSlotRequest updates slot
type HappeningSlotRequest struct {
	HappeningRequest
	SlotID string                        `json:"slotID"`
	Slot   dbo4calendarium.HappeningSlot `json:"slot"`
}

// Validate returns error if not valid
func (v HappeningSlotRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.SlotID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("slotID")

	}
	if err := v.Slot.Validate(); err != nil {
		return err
	}
	return nil
}
