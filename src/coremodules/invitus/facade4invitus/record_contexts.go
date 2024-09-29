package facade4invitus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/invitus/dbo4invitus"
)

// InvitesCollection table name
const InvitesCollection = "invites"

type PersonalInviteEntry = record.DataWithID[string, *dbo4invitus.PersonalInviteDbo]

func NewPersonalInviteEntry(id string) (invite PersonalInviteEntry) {
	return NewPersonalInviteEntryWithDto(id, new(dbo4invitus.PersonalInviteDbo))
}

func NewPersonalInviteEntryWithDto(id string, dbo *dbo4invitus.PersonalInviteDbo) (invite PersonalInviteEntry) {
	invite.ID = id
	invite.Key = NewInviteKey(id)
	invite.Data = dbo
	invite.Record = dal.NewRecordWithData(invite.Key, invite.Data)
	return
}
