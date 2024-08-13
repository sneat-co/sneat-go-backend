package models4auth

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
	"reflect"
)

const PasswordResetKind = "PwdRst"

type PasswordReset struct {
	record.WithID[int]
	Data *PasswordResetData
}

//var _ db.EntityHolder = (*PasswordReset)(nil)

type PasswordResetData struct {
	Email  string
	Status string
	appuser.OwnedByUserWithID
}

func NewPasswordReset(id int, data *PasswordResetData) PasswordReset {
	var key *dal.Key
	if id == 0 {
		key = NewPasswordResetIncompleteKey()
	} else {
		key = NewPasswordResetKey(id)
	}
	if data == nil {
		data = new(PasswordResetData)
	}
	return PasswordReset{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}

func NewPasswordResetKey(id int) *dal.Key {
	return dal.NewKeyWithID(PasswordResetKind, id)
}

func NewPasswordResetIncompleteKey() *dal.Key {
	return dal.NewIncompleteKey(PasswordResetKind, reflect.Int, nil)
}

//func (PasswordReset) Kind() string {
//	return PasswordResetKind
//}
//
////func (record PasswordReset) IntID() int64 {
////	return record.ContactID
////}
//
//func (record PasswordReset) Entity() interface{} {
//	return record.PasswordResetData
//}
//
//func (PasswordReset) NewEntity() interface{} {
//	return new(PasswordResetData)
//}
//
//func (record *PasswordReset) SetEntity(entity interface{}) {
//	if entity == nil {
//		record.PasswordResetData = nil
//	} else {
//		record.PasswordResetData = entity.(*PasswordResetData)
//	}
//}

//func (entity *PasswordResetData) Save() (properties []datastore.Property, err error) {
//	if properties, err = datastore.SaveStruct(entity); err != nil {
//		return
//	}
//	return gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
//		"DtUpdated": gaedb.IsZeroTime,
//		"Email":     gaedb.IsEmptyString,
//	})
//}
