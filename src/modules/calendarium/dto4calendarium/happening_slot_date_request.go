package dto4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strings"
)

// HappeningSlotDateRequest updates slot
type HappeningSlotDateRequest struct {
	HappeningRequest
	Slot models4calendarium.HappeningSlot `json:"slot"`
	Date string                           `json:"date"`
}

// Validate returns error if not valid
func (v HappeningSlotDateRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if err := v.Slot.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.Date) == "" {
		return validation.NewErrRecordIsMissingRequiredField("date")
	}
	if _, err := validate.DateString(v.Date); err != nil {
		return validation.NewErrBadRequestFieldValue("date", err.Error())
	}
	return nil
}
