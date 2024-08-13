package models4auth

import (
	"encoding/base64"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
)

func NewUserEmailKey(email string) *dal.Key {
	return dal.NewKeyWithID(UserEmailKind, GetEmailID(email))
}

func NewUserEmail(email string, data *UserEmailDbo) UserEmailEntry {
	id := GetEmailID(email)
	if data == nil {
		data = new(UserEmailDbo)
	}
	return UserEmailEntry{
		DataWithID: record.NewDataWithID(id, NewUserEmailKey(email), data),
	}
}

var _ appuser.AccountData = (*UserEmailDbo)(nil)

func (entity *UserEmailDbo) GetNames() person.NameFields {
	return entity.NameFields
}

func (entity *UserEmailDbo) ConfirmationPin() string {
	pin := base64.RawURLEncoding.EncodeToString(entity.PasswordBcryptHash)
	//if len(pin) > 20 {
	//	pin = pin[:20]
	//}
	return pin
}

func (entity *UserEmailDbo) IsEmailConfirmed() bool {
	return entity.IsConfirmed
}
