package dto4calendarius

import (
	"fmt"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/strongo/validation"
)

// CreateHappeningRequest DTO
type CreateHappeningRequest struct {
	dto4spaceus.SpaceRequest
	Happening *dbo4calendarius.HappeningBrief `json:"happening"`
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
		return validation.NewErrBadRequestFieldValue("happening", err.Error())
	}
	return nil
}

// CreateHappeningResponse DTO
type CreateHappeningResponse struct {
	ID  string                       `json:"id"`
	Dbo dbo4calendarius.HappeningDbo `json:"dbo"`
}

// Validate returns error if not valid
func (v CreateHappeningResponse) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if err := v.Dbo.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("dbo", err.Error())
	}
	return nil
}
