package facade4listus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
)

// SetListItemsIsDone marks list item as completed
func SetListItemsIsDone(ctx context.Context, userCtx facade.UserContext, request dto4listus.ListItemsSetIsDoneRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4listus.RunListWorker(ctx, userCtx, request.ListRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4listus.ListWorkerParams) (err error) {
			if err = params.GetRecords(ctx, tx); err != nil {
				return
			}
			changed := 0
			for _, item := range params.List.Data.Items {
				for _, id := range request.ItemIDs {
					//logus.Debugf(ctx, "items[%d].InviteID: %s, %s, %s", i, item.ID, "requestItemID", id)
					if item.ID == id {
						if request.IsDone {
							item.Status = const4listus.ListItemStatusDone
						} else if item.IsDone() {
							item.Status = ""
						}
						changed++
					}
				}
			}
			logus.Debugf(ctx, "Number of changed items: %d", changed)
			if changed == 0 {
				return nil
			}
			params.ListUpdates = []dal.Update{
				{
					Field: "items",
					Value: params.List.Data.Items,
				},
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return
}
