package models

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

func NewUserEmail(email string, data *UserEmailData) UserEmail {
	id := GetEmailID(email)
	if data == nil {
		data = new(UserEmailData)
	}
	return UserEmail{
		WithID: record.NewWithID(id, NewUserEmailKey(email), data),
		Data:   data,
	}
}

var _ appuser.AccountData = (*UserEmailData)(nil)

func (entity *UserEmailData) GetNames() person.NameFields {
	return entity.NameFields
}

func (entity *UserEmailData) ConfirmationPin() string {
	pin := base64.RawURLEncoding.EncodeToString(entity.PasswordBcryptHash)
	//if len(pin) > 20 {
	//	pin = pin[:20]
	//}
	return pin
}

func (entity *UserEmailData) IsEmailConfirmed() bool {
	return entity.IsConfirmed
}
