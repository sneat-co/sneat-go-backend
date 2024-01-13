package dto4teamus

import (
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
)

// GetTeamRequest request
type GetTeamRequest struct {
	ID string
}

// Validate validates
func (v *GetTeamRequest) Validate() error {
	id := strings.TrimSpace(v.ID)
	if id == "" {
		return validation.NewErrRecordIsMissingRequiredField("ID")
	}
	if id != v.ID {
		return validation.NewErrBadRequestFieldValue("ID", "has spaces")
	}
	return nil
}

var _ facade.Request = (*GetTeamRequest)(nil)
