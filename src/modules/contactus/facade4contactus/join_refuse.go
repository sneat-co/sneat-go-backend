package facade4contactus

import (
	"context"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// RefuseToJoinSpaceRequest request
type RefuseToJoinSpaceRequest struct {
	SpaceID string `json:"id"`
	Pin     int32  `json:"pin"`
}

// Validate validates request
func (v *RefuseToJoinSpaceRequest) Validate() error {
	if v.SpaceID == "" {
		return validation.NewErrRecordIsMissingRequiredField("space")
	}
	if v.SpaceID == "" {
		return validation.NewErrRecordIsMissingRequiredField("pin")
	}
	return nil
}

// RefuseToJoinSpace refuses to join team
func RefuseToJoinSpace(_ context.Context, userCtx facade.UserContext, request RefuseToJoinSpaceRequest) (err error) {
	err = request.Validate()
	return
}
