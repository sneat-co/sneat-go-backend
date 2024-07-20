package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"time"
)

type InviteRequest struct {
	dto4spaceus.SpaceRequest
	InviteID string `json:"inviteID"`
	Pin      string `json:"pin"`
}

// RejectPersonalInviteRequest holds parameters for rejectio of personal invite
type RejectPersonalInviteRequest = InviteRequest

// Validate validates request
func (v *InviteRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if v.InviteID == "" {
		return validation.NewErrRequestIsMissingRequiredField("inviteID")
	}
	if v.Pin == "" {
		return validation.NewErrRequestIsMissingRequiredField("Pin")
	}
	return nil
}

// RejectPersonalInvite rejects personal invite
func RejectPersonalInvite(ctx context.Context, userContext facade.User, request RejectPersonalInviteRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	team := dal4spaceus.NewSpaceEntry(request.SpaceID)
	invite := NewPersonalInviteEntry(request.InviteID)
	uid := userContext.GetID()

	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		records := []dal.Record{team.Record, invite.Record}
		err := tx.GetMulti(ctx, records)
		if err != nil {
			return err
		}
		if invite.Data.Pin != request.Pin {
			return fmt.Errorf("%w: pin code does not match", facade.ErrBadRequest)
		}
		now := time.Now()
		if err = updateInviteRecord(ctx, tx, uid, now, invite, "rejected"); err != nil {
			return fmt.Errorf("failed to update invite record with rejected status: %w", err)
		}
		return nil
	})
}
