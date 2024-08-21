package models4auth

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"reflect"
	"time"

	"errors"
)

const EmailKind = "Email"

type Email struct {
	record.WithID[int64]
	Data *EmailData
}

func NewEmailKey(id int64) *dal.Key {
	if id == 0 {
		return dal.NewIncompleteKey(EmailKind, reflect.Int64, nil)
	}
	return dal.NewKeyWithID(EmailKind, id)
}

func NewEmail(id int64, data *EmailData) Email {
	key := NewEmailKey(id)
	if data == nil {
		data = new(EmailData)
	}
	return Email{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}

type EmailData struct { // TODO: move into dedicated package
	Status          string
	Error           string `firestore:",omitempty"`
	DtCreated       time.Time
	DtSent          time.Time
	Subject         string `firestore:",omitempty"`
	From            string `firestore:",omitempty"`
	To              string
	BodyText        string `firestore:",omitempty"`
	BodyHtml        string `firestore:",omitempty"`
	AwsSesMessageID string
}

//func (entity *EmailData) Load(ps []datastore.Property) error {
//	return datastore.LoadStruct(entity, ps)
//}

func (entity *EmailData) Validate() (err error) {
	if entity.Status == "" {
		err = errors.New("email.Status is empty")
		return
	}
	if entity.Subject == "" {
		err = errors.New("email.Subject is empty")
		return
	}
	if entity.From == "" {
		err = errors.New("email.From is empty")
		return
	}
	if entity.To == "" {
		err = errors.New("email.To is empty")
		return
	}
	if entity.DtCreated.IsZero() {
		entity.DtCreated = time.Now()
	}
	//if properties, err = datastore.SaveStruct(entity); err != nil {
	//	return
	//}
	//return gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
	//	"DtSent":          gaedb.IsZeroTime,
	//	"AwsSesMessageID": gaedb.IsEmptyString,
	//	"Error":           gaedb.IsEmptyString,
	//	"BodyText":        gaedb.IsEmptyString,
	//	"BodyHtml":        gaedb.IsEmptyString,
	//})
	//return nil, nil
	return nil
}
