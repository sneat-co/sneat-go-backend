package dto4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/strongo/validation"
	"strings"
)

// HappeningRequest DTO
type HappeningRequest struct {
	dto4teamus.SpaceRequest
	HappeningID   string `json:"happeningID"`
	HappeningType string `json:"happeningType,omitempty"`
	//ListType      dbo4listus.ListType `json:"listType,omitempty"` // TODO: Document what it is and why we need it here
}

// Validate returns error if not valid
func (v HappeningRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.HappeningID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("happeningID")
	}
	//switch v.ListType {
	//case "", "to-buy", "to-do":
	//default:
	//	return fmt.Errorf("\"unknown list type: %v", v.ListType)
	//}
	switch v.HappeningType {
	case "", dbo4calendarium.HappeningTypeSingle, dbo4calendarium.HappeningTypeRecurring: // OK
	default:
		return validation.NewErrBadRequestFieldValue("happeningType", "unknown value: "+v.HappeningType)
	}
	return nil
}
