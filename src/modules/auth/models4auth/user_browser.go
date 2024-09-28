package models4auth

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"reflect"
	"time"
)

const UserBrowserKind = "UserBrowser"

type UserBrowserData struct {
	UserID      string
	UserAgent   string
	LastUpdated time.Time `firestore:",omitempty"`
}

type UserBrowser struct {
	record.WithID[int]
	Data *UserBrowserData
}

func NewUserBrowserRecord() dal.Record {
	return dal.NewRecordWithIncompleteKey(UserBrowserKind, reflect.Int, new(UserBrowserData))
}

func NewUserBrowserWithIncompleteKey(data *UserBrowserData) UserBrowser {
	return UserBrowser{
		WithID: record.NewWithID[int](0, dal.NewIncompleteKey(UserBrowserKind, reflect.Int, nil), data),
		Data:   data,
	}
}
