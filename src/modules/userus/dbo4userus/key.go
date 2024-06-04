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

// UserContext defines implementation of `interface facade.UserContext`
type UserContext struct {
	Dbo *UserDbo
	record.WithID[string]
}

// GetID returns user's ID
func (v UserContext) GetID() string {
	return v.WithID.ID
}

// NewUserContext creates new user context
func NewUserContext(id string) (user UserContext) {
	return NewUserContextWithDto(id, new(UserDbo))
}

// NewUserContextWithDto creates new user context with user DTO
func NewUserContextWithDto(id string, dto *UserDbo) (user UserContext) {
	user.WithID.ID = id
	user.Dbo = dto
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
