package facade4splitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-mod-debtus-go/debtus/gae_app/debtstracker/dtdal"
)

var _ dtdal.BillScheduleDal = (*BillScheduleDalGae)(nil)

type BillScheduleDalGae struct {
}

func NewBillScheduleDalGae() BillScheduleDalGae {
	return BillScheduleDalGae{}
}

func (BillScheduleDalGae) GetBillScheduleByID(ctx context.Context, id int64) (models4splitus.BillSchedule, error) {
	key := models4splitus.NewBillScheduleKey(id)
	data := new(models4splitus.BillScheduleEntity)
	billSchedule := models4splitus.BillSchedule{
		WithID: record.WithID[int64]{
			ID:     id,
			Key:    key,
			Record: dal.NewRecordWithData(key, data),
		},
		Data: data,
	}
	db, err := facade.GetSneatDB(ctx)
	if err != nil {
		return billSchedule, err
	}
	if err = db.Get(ctx, billSchedule.Record); err != nil {
		return billSchedule, err
	}
	return billSchedule, err
}

func (BillScheduleDalGae) InsertBillSchedule(_ context.Context, billScheduleEntity *models4splitus.BillScheduleEntity) (billSchedule models4splitus.BillSchedule, err error) {
	_ = models4splitus.NewBillScheduleIncompleteKey()
	panic("TODO: implement me")
	//key := NewBillScheduleIncompleteKey()
	//if key, err = gaedb.Put(ctx, key, billScheduleEntity); err != nil {
	//	return
	//}
	//billSchedule.ContactID = key.ContactID.(int)
	//return
}

func (BillScheduleDalGae) UpdateBillSchedule(ctx context.Context, billSchedule models4splitus.BillSchedule) error {
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(ctx, billSchedule.Record)
	})
}
