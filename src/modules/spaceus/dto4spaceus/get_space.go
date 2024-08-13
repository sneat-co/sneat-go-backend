package dto4spaceus

import (
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
)

// GetSpaceRequest request
type GetSpaceRequest struct {
	ID string
}

// Validate validates
func (v *GetSpaceRequest) Validate() error {
	id := strings.TrimSpace(v.ID)
	if id == "" {
		return validation.NewErrRecordIsMissingRequiredField("ContactID")
	}
	if id != v.ID {
		return validation.NewErrBadRequestFieldValue("ContactID", "has spaces")
	}
	return nil
}

var _ facade.Request = (*GetSpaceRequest)(nil)
