package dto4teamus

import (
	"github.com/strongo/validation"
	"strings"
)

// NewTeamRequest creates new team request
func NewTeamRequest(teamID string) TeamRequest {
	return TeamRequest{TeamID: teamID}
}

// TeamRequest request
type TeamRequest struct {
	TeamID string `json:"teamID"`
}

// Validate validates request
func (v *TeamRequest) Validate() error {
	if strings.TrimSpace(v.TeamID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("team")
	}
	return nil
}
