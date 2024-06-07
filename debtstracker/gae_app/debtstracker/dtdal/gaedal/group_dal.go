package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/delaying"
	"github.com/strongo/log"
)

var _ dtdal.GroupDal = (*GroupDalGae)(nil)

type GroupDalGae struct { // TODO: Obsolete naming with migration to Dalgo
}

func NewGroupDalGae() GroupDalGae {
	return GroupDalGae{}
}

func (GroupDalGae) InsertGroup(c context.Context, tx dal.ReadwriteTransaction, groupEntity *models.GroupDbo) (group models.GroupEntry, err error) {
	group = models.NewGroup("", groupEntity)
	err = dtdal.InsertWithRandomStringID(c, tx, group.Record)
	return
}

func (GroupDalGae) SaveGroup(c context.Context, tx dal.ReadwriteTransaction, group models.GroupEntry) (err error) {
	if err = tx.Set(c, group.Record); err != nil {
		return
	}
	return
}

func (GroupDalGae) GetGroupByID(c context.Context, tx dal.ReadSession, groupID string) (group models.GroupEntry, err error) {
	if tx == nil {
		if tx, err = facade.GetDatabase(c); err != nil {
			return
		}
	}
	if group.ID = groupID; group.ID == "" {
		panic("groupID is empty string")
	}
	group = models.NewGroup(groupID, nil)
	if err = tx.Get(c, group.Record); err != nil {
		return
	}
	return
}

func (GroupDalGae) DelayUpdateGroupWithBill(c context.Context, groupID, billID string) (err error) {
	if err = delayUpdateGroupWithBill.EnqueueWork(c, delaying.With(common.QUEUE_BILLS, "UpdateGroupWithBill", 0), groupID, billID); err != nil {
		return
	}
	return
}

func delayedUpdateGroupWithBill(c context.Context, groupID, billID string) (err error) {
	log.Debugf(c, "delayedUpdateGroupWithBill(groupID=%d, billID=%d)", groupID, billID)
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		bill, err := facade.GetBillByID(c, tx, billID)
		if err != nil {
			return
		}
		var group models.GroupEntry
		if group, err = dtdal.Group.GetGroupByID(c, tx, groupID); err != nil {
			return err
		}
		var changed bool
		if changed, err = group.Data.AddBill(bill); err != nil {
			return err
		} else if changed {
			if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
				return err
			}
		}
		return
	}); err != nil {
		return
	}
	return
}
