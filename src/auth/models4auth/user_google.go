package models4auth

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
)

const UserGoogleCollection = "UserAccountEntry"
const UserFacebookCollection = "UserFb"

// UserAccountEntry - TODO: consider migrating to https://github.com/dal-go/dalgo4auth
type UserAccountEntry = record.DataWithID[string, *appuser.AccountDataBase]

func GetUserAccountKey(ua UserAccountEntry) appuser.AccountKey {
	return ua.Data.AccountKey
}

func NewUserAccountEntry(id string) UserAccountEntry {
	key := dal.NewKeyWithID(UserGoogleCollection, id)
	data := new(appuser.AccountDataBase)
	return record.NewDataWithID(id, key, data)
}
