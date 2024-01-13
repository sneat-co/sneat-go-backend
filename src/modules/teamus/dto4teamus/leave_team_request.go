package dto4teamus

import (
	"github.com/strongo/validation"
	"strings"
)

// LeaveTeamRequest request
type LeaveTeamRequest struct {
	TeamRequest
	Message string `json:"message,omitempty"`
}

// Validate validates request
func (v *LeaveTeamRequest) Validate() error {
	if v.TeamID = strings.TrimSpace(v.TeamID); v.TeamID == "" {
		return validation.NewErrRecordIsMissingRequiredField("TeamIDs")
	}
	return nil
}
