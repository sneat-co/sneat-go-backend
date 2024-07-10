package dto4teamus

import (
	"github.com/strongo/validation"
	"strings"
)

// NewSpaceRequest creates new team request
func NewSpaceRequest(teamID string) SpaceRequest {
	return SpaceRequest{SpaceID: teamID}
}

// SpaceRequest request
type SpaceRequest struct {
	SpaceID string `json:"spaceID"`
}

// Validate validates request
func (v *SpaceRequest) Validate() error {
	if strings.TrimSpace(v.SpaceID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("space")
	}
	return nil
}
