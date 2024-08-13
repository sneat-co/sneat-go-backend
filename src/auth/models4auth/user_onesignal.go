package models4auth

import (
	"github.com/dal-go/dalgo/record"
	"time"
)

const UserOneSignalKind = "UserOneSignal"

type UserOneSignalEntity struct {
	UserID  int64
	Created time.Time
}

type UserOneSignal struct {
	record.WithID[string]
	*UserOneSignalEntity
}

//var _ db.EntityHolder = (*UserOneSignal)(nil)

func (UserOneSignal) Kind() string {
	return UserOneSignalKind
}

func (userOneSignal UserOneSignal) Entity() interface{} {
	return userOneSignal.UserOneSignalEntity
}

func (UserOneSignal) NewEntity() interface{} {
	return new(UserOneSignalEntity)
}

func (userOneSignal *UserOneSignal) SetEntity(entity interface{}) {
	userOneSignal.UserOneSignalEntity = entity.(*UserOneSignalEntity)
}
