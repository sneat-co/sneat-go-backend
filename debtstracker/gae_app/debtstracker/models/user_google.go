package models

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
)

const UserGoogleCollection = "UserAccount"
const UserFacebookCollection = "UserFb"

var _ appuser.AccountRecord = (*UserAccount)(nil)

// UserAccount - TODO: consider migrating to https://github.com/dal-go/dalgo4auth
type UserAccount struct { // TODO: Move out to library?
	record.WithID[string]
	data *appuser.AccountDataBase
}

func (ua UserAccount) Key() appuser.AccountKey {
	return ua.data.AccountKey
}

func (ua UserAccount) Data() appuser.AccountData {
	return ua.data
}

func (ua UserAccount) DataStruct() *appuser.AccountDataBase {
	return ua.data
}

func NewUserAccount(id string) UserAccount {
	key := dal.NewKeyWithID(UserGoogleCollection, id)
	data := new(appuser.AccountDataBase)
	return UserAccount{
		WithID: record.WithID[string]{
			ID:     id,
			Key:    key,
			Record: dal.NewRecordWithData(key, data),
		},
		data: data,
	}
}
