package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

var _ dtdal.BillScheduleDal = (*BillScheduleDalGae)(nil)

type BillScheduleDalGae struct {
}

func NewBillScheduleDalGae() BillScheduleDalGae {
	return BillScheduleDalGae{}
}

func (BillScheduleDalGae) GetBillScheduleByID(c context.Context, id int64) (models.BillSchedule, error) {
	key := models.NewBillScheduleKey(id)
	data := new(models.BillScheduleEntity)
	billSchedule := models.BillSchedule{
		WithID: record.WithID[int64]{
			ID:     id,
			Key:    key,
			Record: dal.NewRecordWithData(key, data),
		},
		Data: data,
	}
	db, err := GetDatabase(c)
	if err != nil {
		return billSchedule, err
	}
	if err = db.Get(c, billSchedule.Record); err != nil {
		return billSchedule, err
	}
	return billSchedule, err
}

func (BillScheduleDalGae) InsertBillSchedule(c context.Context, billScheduleEntity *models.BillScheduleEntity) (billSchedule models.BillSchedule, err error) {
	_ = models.NewBillScheduleIncompleteKey()
	panic("TODO: implement me")
	//key := NewBillScheduleIncompleteKey()
	//if key, err = gaedb.Put(c, key, billScheduleEntity); err != nil {
	//	return
	//}
	//billSchedule.ID = key.ID.(int)
	//return
}

func (BillScheduleDalGae) UpdateBillSchedule(c context.Context, billSchedule models.BillSchedule) error {
	db, err := GetDatabase(c)
	if err != nil {
		return err
	}
	return db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(c, billSchedule.Record)
	})
}
