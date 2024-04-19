package dto4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/strongo/validation"
	"strings"
)

// HappeningRequest DTO
type HappeningRequest struct {
	dto4teamus.TeamRequest
	HappeningID   string `json:"happeningID"`
	HappeningType string `json:"happeningType,omitempty"`
	//ListType      models4listus.ListType `json:"listType,omitempty"` // TODO: Document what it is and why we need it here
}

// Validate returns error if not valid
func (v HappeningRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
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
	case "", "single", "recurring": // OK
	default:
		return validation.NewErrBadRequestFieldValue("happeningType", "unknown value: "+v.HappeningType)
	}
	return nil
}
