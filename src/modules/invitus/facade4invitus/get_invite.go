package facade4invitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/dbo4invitus"
)

// GetInviteByID returns an invitation record by ID
func GetInviteByID(ctx context.Context, getter dal.ReadSession, id string) (inviteDto *dbo4invitus.InviteDbo, inviteRecord dal.Record, err error) {
	inviteDto = new(dbo4invitus.InviteDbo)
	inviteRecord = dal.NewRecordWithData(NewInviteKey(id), inviteDto)
	return inviteDto, inviteRecord, getter.Get(ctx, inviteRecord)
}

// GetPersonalInviteByID returns an invitation record by ID
func GetPersonalInviteByID(ctx context.Context, getter dal.ReadSession, id string) (inviteDto *dbo4invitus.PersonalInviteDbo, inviteRecord dal.Record, err error) {
	inviteDto = new(dbo4invitus.PersonalInviteDbo)
	inviteRecord = dal.NewRecordWithData(NewInviteKey(id), inviteDto)
	return inviteDto, inviteRecord, getter.Get(ctx, inviteRecord)
}
