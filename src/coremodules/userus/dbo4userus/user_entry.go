package dbo4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

// UserEntry wraps UserDbo with dalgo record
type UserEntry struct {
	record.DataWithID[string, *UserDbo]
}

// GetID returns user ID - needed for BotsGoFramework
func (v UserEntry) GetID() string {
	return v.ID
}

// NewUserEntry creates new user context
func NewUserEntry(id string) (user UserEntry) {
	return NewUserEntryWithDbo(id, new(UserDbo))
}

// NewUserEntryWithDbo creates new user context with user DTO
func NewUserEntryWithDbo(id string, dto *UserDbo) (user UserEntry) {
	key := NewUserKey(id)
	user.WithID = record.WithID[string]{
		ID:     id,
		FullID: "users/" + id,
		Key:    key,
	}
	user.Data = dto
	user.Record = dal.NewRecordWithData(user.Key, dto)
	return
}
