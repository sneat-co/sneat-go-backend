package dto4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
)

// HappeningSlotRequest updates slot
type HappeningSlotRequest struct {
	HappeningRequest
	Slot dbo4calendarium.HappeningSlot `json:"slot"`
}

// Validate returns error if not valid
func (v HappeningSlotRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if err := v.Slot.Validate(); err != nil {
		return err
	}
	return nil
}
