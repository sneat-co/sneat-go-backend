package facade4splitus

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/logus"
)

func SaveBill(c context.Context, tx dal.ReadwriteTransaction, bill models4splitus.BillEntry) (err error) {
	if err = tx.Set(c, bill.Record); err != nil {
		return
	}
	if err = DelayUpdateUsersWithBill(c, bill.ID, bill.Data.UserIDs); err != nil {
		return
	}
	return
}

func delayedUpdateBillDependencies(c context.Context, billID string) (err error) {
	logus.Debugf(c, "delayerUpdateBillDependencies(billID=%s)", billID)
	var bill models4splitus.BillEntry
	if bill, err = GetBillByID(c, nil, billID); err != nil {
		if dal.IsNotFound(err) {
			logus.Warningf(c, err.Error())
			err = nil
		}
		return
	}
	if userGroupID := bill.Data.GetUserGroupID(); userGroupID != "" {
		if err = DelayUpdateGroupWithBill(c, userGroupID, bill.ID); err != nil {
			return
		}
	}
	for _, member := range bill.Data.GetBillMembers() {
		if member.UserID != "" {
			if err = DelayUpdateSpaceWithBill(c, member.UserID, bill.ID); err != nil {
				return
			}
		}
	}
	return
}

func UpdateBillsHolder(c context.Context, tx dal.ReadwriteTransaction, billID string, getBillsHolder dtdal.BillsHolderGetter) (err error) {
	_, _, _, _ = c, tx, billID, getBillsHolder
	return errors.New("UpdateBillsHolder() is not implemented yet")
}
