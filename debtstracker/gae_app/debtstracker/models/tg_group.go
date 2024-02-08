package models

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

const TgGroupKind = "TgGroup"

type TgGroup struct {
	record.WithID[int64]
	Data *TgGroupData
}

func NewTgGroup(id int64, data *TgGroupData) TgGroup {
	key := dal.NewKeyWithID(TgGroupKind, id)
	return TgGroup{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}

//var _ db.EntityHolder = (*TgGroup)(nil)

type TgGroupData struct {
	UserGroupID string `datastore:",noindex" firestore:",noindex"`
}

//func (TgGroup) Kind() string {
//	return TgGroupKind
//}
//
//func (tgGroup TgGroup) Entity() interface{} {
//	return tgGroup.TgGroupData
//}
//
//func (tgGroup TgGroup) NewEntity() interface{} {
//	return new(TgGroupData)
//}
//
//func (tgGroup *TgGroup) SetEntity(entity interface{}) {
//	if entity == nil {
//		tgGroup.TgGroupData = nil
//	} else {
//		tgGroup.TgGroupData = entity.(*TgGroupData)
//	}
//}
