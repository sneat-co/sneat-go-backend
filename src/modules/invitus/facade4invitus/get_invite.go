package facade4invitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/models4invitus"
)

// GetInviteByID returns an invitation record by ID
func GetInviteByID(ctx context.Context, getter dal.ReadSession, id string) (inviteDto *models4invitus.InviteDto, inviteRecord dal.Record, err error) {
	inviteDto = new(models4invitus.InviteDto)
	inviteRecord = dal.NewRecordWithData(NewInviteKey(id), inviteDto)
	return inviteDto, inviteRecord, getter.Get(ctx, inviteRecord)
}

// GetPersonalInviteByID returns an invitation record by ID
func GetPersonalInviteByID(ctx context.Context, getter dal.ReadSession, id string) (inviteDto *models4invitus.PersonalInviteDto, inviteRecord dal.Record, err error) {
	inviteDto = new(models4invitus.PersonalInviteDto)
	inviteRecord = dal.NewRecordWithData(NewInviteKey(id), inviteDto)
	return inviteDto, inviteRecord, getter.Get(ctx, inviteRecord)
}
