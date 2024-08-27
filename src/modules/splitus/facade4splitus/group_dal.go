package facade4splitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/dal4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
)

func delayedUpdateGroupWithBill(ctx context.Context, spaceID, billID string) (err error) {
	logus.Debugf(ctx, "delayedUpdateGroupWithBill(spaceID=%s, billID=%s)", spaceID, billID)
	if err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		bill, err := GetBillByID(ctx, tx, billID)
		if err != nil {
			return
		}
		splitusSpace := models4splitus.NewSplitusSpaceEntry(spaceID)
		if err = dal4splitus.GetSplitusSpace(ctx, tx, splitusSpace); err != nil {
			return err
		}
		var changed bool
		if changed, err = splitusSpace.Data.AddBill(bill); err != nil {
			return err
		} else if changed {
			if err = dal4splitus.SaveSplitusSpace(ctx, tx, splitusSpace); err != nil {
				return err
			}
		}
		return
	}); err != nil {
		return
	}
	return
}
