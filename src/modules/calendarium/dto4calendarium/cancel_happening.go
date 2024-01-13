package dto4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strings"
)

type CancelHappeningRequest struct {
	HappeningRequest
	Date   string `json:"date,omitempty"`
	SlotID string `json:"slotID,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func (v CancelHappeningRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if v.Date != "" {
		if _, err := validate.DateString(v.Date); err != nil {
			return validation.NewErrBadRequestFieldValue("date", err.Error())
		}
	}
	if v.Date != "" && strings.TrimSpace(v.SlotID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("slotIDs")
	}
	if len(v.Reason) > models4calendarium.ReasonMaxLen {
		return validation.NewErrBadRequestFieldValue("reason",
			fmt.Sprintf("maximum length of reason is %v, got %v", models4calendarium.ReasonMaxLen, len(v.Reason)))
	}
	return nil
}
