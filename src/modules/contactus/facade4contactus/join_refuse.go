package facade4contactus

import (
	"context"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// RefuseToJoinTeamRequest request
type RefuseToJoinTeamRequest struct {
	TeamID string `json:"id"`
	Pin    int32  `json:"pin"`
}

// Validate validates request
func (v *RefuseToJoinTeamRequest) Validate() error {
	if v.TeamID == "" {
		return validation.NewErrRecordIsMissingRequiredField("team")
	}
	if v.TeamID == "" {
		return validation.NewErrRecordIsMissingRequiredField("pin")
	}
	return nil
}

// RefuseToJoinTeam refuses to join team
func RefuseToJoinTeam(_ context.Context, userContext facade.User, request RefuseToJoinTeamRequest) (err error) {
	err = request.Validate()
	return
}
