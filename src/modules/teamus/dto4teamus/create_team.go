package dto4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
)

var _ facade.Request = (*CreateTeamRequest)(nil)

// CreateTeamRequest request
type CreateTeamRequest struct {
	Type  core4teamus.TeamType `json:"type"`
	Title string               `json:"title,omitempty"`
}

// Validate validates request
func (request *CreateTeamRequest) Validate() error {
	if strings.TrimSpace(string(request.Type)) == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if request.Type != "family" && strings.TrimSpace(request.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// CreateTeamResponse response
type CreateTeamResponse = TeamResponse
