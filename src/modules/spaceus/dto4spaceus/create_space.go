package dto4spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
)

var _ facade.Request = (*CreateSpaceRequest)(nil)

// CreateSpaceRequest request
type CreateSpaceRequest struct {
	Type  core4teamus.SpaceType `json:"type"`
	Title string                `json:"title,omitempty"`
}

// Validate validates request
func (request *CreateSpaceRequest) Validate() error {
	if strings.TrimSpace(string(request.Type)) == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if request.Type != "family" && strings.TrimSpace(request.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// CreateSpaceResponse response
type CreateSpaceResponse = SpaceResponse
