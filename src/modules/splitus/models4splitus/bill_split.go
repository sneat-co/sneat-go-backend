package models4splitus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
)

type Split struct {
	record.WithID[int]
	*SplitEntity
}

//var _ db.EntityHolder = (*Split)(nil)

type SplitEntity struct {
	appuser.OwnedByUserWithID
	BillIDs []string `firestore:",omitempty"`
}

func (Split) Kind() string {
	return SplitsCollection
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
