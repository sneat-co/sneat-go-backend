package models4splitus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-mod-debtus-go/debtus/const4debtus"
	"reflect"
	"time"
)

type BillScheduleStatus string

const (
	BillScheduleStatusDraft    BillScheduleStatus = "draft"
	BillScheduleStatusActive   BillScheduleStatus = const4debtus.StatusActive
	BillScheduleStatusArchived BillScheduleStatus = const4debtus.StatusArchived
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
	BillsCount        int    `firestore:",omitempty"`
	CreatedFromBillID string `firestore:",omitempty"`
	RepeatPeriod      Period `firestore:",omitempty"`
	RepeatOn          string `firestore:",omitempty"`
	IsAutoTransfer    bool   `firestore:",omitempty"`

	LastBillID string    `firestore:",omitempty"`
	DtLast     time.Time `firestore:",omitempty"`
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
//	return bill.ContactID
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
