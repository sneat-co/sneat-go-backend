package models

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"reflect"
	"time"
)

type BillScheduleStatus string

const (
	BillScheduleStatusDraft    BillScheduleStatus = "draft"
	BillScheduleStatusActive   BillScheduleStatus = STATUS_ACTIVE
	BillScheduleStatusArchived BillScheduleStatus = STATUS_ARCHIVED
	//BillScheduleStatusDeleted  BillScheduleStatus = STATUS_DELETED
)

type Period string

const (
	PeriodWeekly  Period = "weekly"
	PeriodMonthly Period = "monthly"
	PeriodYearly  Period = "yearly"
)

const BillScheduleKind = "BillSchedule"

type BillSchedule struct {
	record.WithID[int64]
	Data *BillScheduleEntity
}

func NewBillScheduleKey(id int64) *dal.Key {
	return dal.NewKeyWithID(BillScheduleKind, id)
}

func NewBillScheduleIncompleteKey() *dal.Key {
	return dal.NewIncompleteKey(BillScheduleKind, reflect.Int64, nil)
}

type BillScheduleEntity struct {
	BillCommon
	/* Repeat examples (RepeatPeriod:RepeatOn)
	* weekly:monday
	* monthly:2 - 2nd day of each month. possible values 1-28
	// * monthly:first-monday
	// * yearly:1-jan ???
	*/
	BillsCount        int    `datastore:",noindex"`
	CreatedFromBillID string `datastore:",noindex"`
	RepeatPeriod      Period `datastore:",noindex"`
	RepeatOn          string `datastore:",noindex"`
	IsAutoTransfer    bool   `datastore:",noindex"`

	LastBillID string    `datastore:",noindex"`
	DtLast     time.Time `datastore:",noindex"`
	DtNext     time.Time
}

func (entity *BillScheduleEntity) Validate() (err error) {
	//if properties, err = datastore.SaveStruct(entity); err != nil {
	//	return
	//}
	if err = entity.BillCommon.Validate(); err != nil {
		return
	}
	//if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
	//	"DtLast":            gaedb.IsZeroTime,
	//	"DtNext":            gaedb.IsZeroTime,
	//	"LastBillID":        gaedb.IsZeroInt,
	//	"IsAutoTransfer":    gaedb.IsZeroBool,
	//	"BillsCount":        gaedb.IsZeroInt,
	//	"CreatedFromBillID": gaedb.IsZeroInt,
	//}); err != nil {
	//	return
	//}
	return
}

func (BillSchedule) Kind() string {
	return BillKind
}

//func (bill BillSchedule) IntID() int64 {
//	return bill.ID
//}

func (bill *BillSchedule) Entity() interface{} {
	if bill.Data == nil {
		bill.Data = new(BillScheduleEntity)
	}
	return bill.Data
}

func (bill *BillSchedule) SetEntity(entity interface{}) {
	bill.Data = entity.(*BillScheduleEntity)
}
