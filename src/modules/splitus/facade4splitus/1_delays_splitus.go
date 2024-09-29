package facade4splitus

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/userus/const4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/const4splitus"
	"github.com/strongo/delaying"
)

var delayerUpdateUserWithGroups delaying.Function

func delayUpdateUserWithGroups(ctx context.Context, userID string, groupIDs2add, groupIDs2remove []string) (err error) { // TODO: make args meaningful
	args := []any{userID, groupIDs2add, groupIDs2remove}
	params := delaying.With(const4userus.QueueUsers, "update-user-with-groups", 0)
	return delayerUpdateUserWithGroups.EnqueueWorkMulti(ctx, params, args)
}

var delayerUpdateGroupUsers delaying.Function
var delayerUpdateContactWithGroups delaying.Function

// ------------------------------------------------------------
var delayerUpdateGroupWithBill delaying.Function

func DelayUpdateGroupWithBill(ctx context.Context, groupID, billID string) (err error) {
	if err = delayerUpdateGroupWithBill.EnqueueWork(ctx, delaying.With(const4splitus.QueueSplitus, "UpdateGroupWithBill", 0), groupID, billID); err != nil {
		return
	}
	return
}

// ------------------------------------------------------------
var delayerUpdateBillDependencies delaying.Function

func DelayUpdateBillDependencies(ctx context.Context, billID string) (err error) {
	if err = delayerUpdateBillDependencies.EnqueueWork(ctx, delaying.With(const4splitus.QueueSplitus, "UpdateBillDependencies", 0), billID); err != nil {
		return
	}
	return
}

var delayerUpdateUsersWithBill delaying.Function
var delayerUpdateUserWithBill delaying.Function

//------------------------------------------------------------

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Function) {
	delayerUpdateUserWithGroups = mustRegisterFunc("delayedUpdateUserWithGroups", delayedUpdateUserWithGroups)
	delayerUpdateGroupUsers = mustRegisterFunc("delayedUpdateGroupUsers", delayedUpdateGroupUsers)
	delayerUpdateContactWithGroups = mustRegisterFunc("delayedUpdateContactWithGroup", delayedUpdateContactWithGroup)
	delayerUpdateGroupWithBill = mustRegisterFunc("delayedUpdateContactWithGroup", delayedUpdateGroupWithBill)
	delayerUpdateBillDependencies = mustRegisterFunc("delayerUpdateBillDependencies", delayedUpdateBillDependencies)
	delayerUpdateUsersWithBill = mustRegisterFunc(updateUsersWithBillKeyName, delayedUpdateUsersWithBill)
	delayerUpdateGroupWithBill = mustRegisterFunc("delayedUpdateWithBill", delayedUpdateGroupWithBill)
	delayerUpdateUserWithBill = mustRegisterFunc("delayedUpdateUserWithBill", delayedUpdateUserWithBill)

}
