package dbo4userus

import (
	"github.com/dal-go/dalgo/dal"
)

// Kind is defining collection name for users records
const Kind = "users"

// NewUserKey creates new user doc ref
func NewUserKey(id string) *dal.Key {
	return dal.NewKeyWithID(Kind, id)
}

// NewUserKeys creates new api4meetingus doc refs
func NewUserKeys(ids []string) (userKeys []*dal.Key) {
	userKeys = make([]*dal.Key, len(ids))
	for i, id := range ids {
		userKeys[i] = NewUserKey(id)
	}
	return userKeys
}
