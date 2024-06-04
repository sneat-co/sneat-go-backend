package facade4invitus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/dbo4invitus"
)

// InvitesCollection table name
const InvitesCollection = "invites"

type PersonalInviteContext struct {
	ID  string
	Dto *dbo4invitus.PersonalInviteDto
	record.WithID[string]
}

func NewPersonalInviteContext(id string) (invite PersonalInviteContext) {
	return NewPersonalInviteContextWithDto(id, new(dbo4invitus.PersonalInviteDto))
}

func NewPersonalInviteContextWithDto(id string, dto *dbo4invitus.PersonalInviteDto) (invite PersonalInviteContext) {
	invite.ID = id
	invite.Key = NewInviteKey(id)
	invite.Dto = dto
	invite.Record = dal.NewRecordWithData(invite.Key, invite.Dto)
	return
}
