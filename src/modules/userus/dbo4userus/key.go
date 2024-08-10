package dbo4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

// Kind is defining collection name for users records
const Kind = "users"

// NewUserKey creates new user doc ref
func NewUserKey(id string) *dal.Key {
	return dal.NewKeyWithID(Kind, id)
}

// UserEntry defines implementation of `interface facade2debtus.UserEntry`
type UserEntry struct {
	record.DataWithID[string, *UserDbo]
}

func (v UserEntry) GetID() string {
	return v.ID
}

// NewUserEntry creates new user context
func NewUserEntry(id string) (user UserEntry) {
	return NewUserContextWithDto(id, new(UserDbo))
}

// NewUserContextWithDto creates new user context with user DTO
func NewUserContextWithDto(id string, dto *UserDbo) (user UserEntry) {
	user.WithID.ID = id
	user.Data = dto
	user.Key = NewUserKey(id)
	user.Record = dal.NewRecordWithData(user.Key, dto)
	return
}

// NewUserKeys creates new api4meetingus doc refs
func NewUserKeys(ids []string) (userKeys []*dal.Key) {
	userKeys = make([]*dal.Key, len(ids))
	for i, id := range ids {
		userKeys[i] = NewUserKey(id)
	}
	return userKeys
}
