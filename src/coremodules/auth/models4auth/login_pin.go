package models4auth

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"reflect"
	"time"
)

const LoginPinKind = "LoginPin"

// LoginPin - TODO check and describe how it is different from LoginCode
type LoginPin struct {
	record.WithID[int]
	Data *LoginPinData
}

//var _ db.EntityHolder = (*LoginPin)(nil)

//func (LoginPin) Kind() string {
//	return LoginPinKind
//}
//
//func (loginPin LoginPin) Entity() interface{} {
//	return loginPin.LoginPinData
//}
//
//func (LoginPin) NewEntity() interface{} {
//	return new(LoginPinData)
//}
//
//func (loginPin *LoginPin) SetEntity(entity interface{}) {
//	if entity == nil {
//		loginPin.LoginPinData = nil
//	} else {
//		loginPin.LoginPinData = entity.(*LoginPinData)
//	}
//
//}

// LoginPinData is a data structure for LoginPin entity.
// TODO check and describe how it is different from LoginCodeData
type LoginPinData struct {
	Channel    string `firestore:",omitempty"`
	GaClientID string `firestore:",omitempty"`
	Created    time.Time
	Pinned     time.Time `firestore:",omitempty"`
	SignedIn   time.Time `firestore:",omitempty"`
	UserID     string    `firestore:",omitempty"`
	Code       int32     `firestore:",omitempty"`
}

func (entity *LoginPinData) IsActive(channel string) bool {
	return entity.SignedIn.IsZero() && entity.Channel == channel
}

func NewLoginPinIncompleteKey() *dal.Key {
	return dal.NewIncompleteKey(LoginPinKind, reflect.Int, nil)
}

func NewLoginPinKey(id int) *dal.Key {
	if id == 0 {
		return NewLoginPinIncompleteKey()
	}
	return dal.NewKeyWithID(LoginPinKind, id)
}

func NewLoginPin(id int, data *LoginPinData) LoginPin {
	if data == nil {
		data = new(LoginPinData)
	}
	key := NewLoginPinKey(id)
	return LoginPin{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}
