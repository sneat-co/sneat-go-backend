package models

import (
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
)

const SplitKind = "Split"

type Split struct {
	record.WithID[int]
	*SplitEntity
}

//var _ db.EntityHolder = (*Split)(nil)

type SplitEntity struct {
	appuser.OwnedByUserWithID
	BillIDs []string `datastore:",noindex"`
}

func (Split) Kind() string {
	return SplitKind
}

func (record Split) Entity() interface{} {
	return record.SplitEntity
}

func (Split) NewEntity() interface{} {
	return new(SplitEntity)
}

func (record *Split) SetEntity(entity interface{}) {
	if entity == nil {
		record.SplitEntity = nil
	} else {
		record.SplitEntity = entity.(*SplitEntity)
	}

}
