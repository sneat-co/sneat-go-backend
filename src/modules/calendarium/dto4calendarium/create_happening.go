package dto4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/strongo/validation"
)

// CreateHappeningRequest DTO
type CreateHappeningRequest struct {
	dto4spaceus.SpaceRequest
	Happening *dbo4calendarium.HappeningBrief `json:"happening"`
}

// Validate returns error if not valid
func (v CreateHappeningRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return fmt.Errorf("space request is not valid: %w", err)
	}
	if v.Happening == nil {
		return validation.NewErrRequestIsMissingRequiredField("happening")
	}
	if err := v.Happening.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("dto4debtus", err.Error())
	}
	return nil
}

// CreateHappeningResponse DTO
type CreateHappeningResponse struct {
	ID  string                       `json:"id"`
	Dto dbo4calendarium.HappeningDbo `json:"dto4debtus"`
}

// Validate returns error if not valid
func (v CreateHappeningResponse) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	return nil
}
