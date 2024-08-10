package gaedal

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
)

type billDalGae struct {
}

var _ dtdal.BillDal = (*billDalGae)(nil) // Make sure we implement interface

func newBillDalGae() billDalGae {
	return billDalGae{}
}

func (billDalGae) SaveBill(c context.Context, tx dal.ReadwriteTransaction, bill models.Bill) (err error) {
	if err = tx.Set(c, bill.Record); err != nil {
		return
	}
	if err = DelayUpdateUsersWithBill(c, bill.ID, bill.Data.UserIDs); err != nil {
		return
	}
	return
}

func (billDalGae) DelayUpdateBillDependencies(c context.Context, billID string) (err error) {
	if err = delayUpdateBillDependencies.EnqueueWork(c, delaying.With(common.QUEUE_BILLS, "UpdateBillDependencies", 0), billID); err != nil {
		return
	}
	return
}

func delayedUpdateBillDependencies(c context.Context, billID string) (err error) {
	logus.Debugf(c, "delayUpdateBillDependencies(billID=%s)", billID)
	var bill models.Bill
	if bill, err = facade2debtus.GetBillByID(c, nil, billID); err != nil {
		if dal.IsNotFound(err) {
			logus.Warningf(c, err.Error())
			err = nil
		}
		return
	}
	if userGroupID := bill.Data.GetUserGroupID(); userGroupID != "" {
		if err = dtdal.Group.DelayUpdateGroupWithBill(c, userGroupID, bill.ID); err != nil {
			return
		}
	}
	for _, member := range bill.Data.GetBillMembers() {
		if member.UserID != "" {
			if err = dtdal.User.DelayUpdateUserWithBill(c, member.UserID, bill.ID); err != nil {
				return
			}
		}
	}
	return
}

func (billDalGae) UpdateBillsHolder(c context.Context, tx dal.ReadwriteTransaction, billID string, getBillsHolder dtdal.BillsHolderGetter) (err error) {
	_, _, _, _ = c, tx, billID, getBillsHolder
	return errors.New("UpdateBillsHolder() is not implemented yet")
}
