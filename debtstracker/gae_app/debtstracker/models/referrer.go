package models

import (
	"github.com/dal-go/dalgo/record"
	"time"
)

const RefererKind = "Referer"

type Referer struct {
	record.WithID[int]
	Data *RefererEntity
}

//var _ db.EntityHolder = (*Referer)(nil)

//func (Referer) Kind() string {
//	return RefererKind
//}

//func (r Referer) Entity() interface{} {
//	return r.RefererEntity
//}
//
//func (Referer) NewEntity() interface{} {
//	return new(RefererEntity)
//}
//
//func (r *Referer) SetEntity(entity interface{}) {
//	if entity == nil {
//		r.RefererEntity = nil
//	} else {
//		r.RefererEntity = entity.(*RefererEntity)
//	}
//}

type RefererEntity struct {
	Platform   string    `datastore:"p"`
	ReferredTo string    `datastore:"to"`
	ReferredBy string    `datastore:"by"`
	DtCreated  time.Time `datastore:"t"`
}
