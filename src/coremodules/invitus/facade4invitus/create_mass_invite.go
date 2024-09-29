package facade4invitus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/invitus/dbo4invitus"
)

// CreateMassInviteRequest parameters for creating a mass invite
type CreateMassInviteRequest struct {
	Invite dbo4invitus.MassInvite `json:"invite"`
}

// Validate validates parameters for creating a mass invite
func (request *CreateMassInviteRequest) Validate() error {
	return request.Invite.Validate()
}

// CreateMassInviteResponse creating a mass invite
type CreateMassInviteResponse struct {
	ID string `json:"id"`
}

// CreateMassInvite creates a mass invite
func CreateMassInvite(_ context.Context, _ CreateMassInviteRequest) (response CreateMassInviteResponse, err error) {
	//request.InviteDbo.SpaceIDs.InviteID
	response.ID = ""
	return
}
