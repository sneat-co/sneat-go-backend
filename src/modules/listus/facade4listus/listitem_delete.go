package facade4listus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteListItems deletes list items
func DeleteListItems(ctx context.Context, userCtx facade.UserContext, request dto4listus.ListItemIDsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4listus.RunListWorker(ctx, userCtx, request.ListRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4listus.ListWorkerParams) error {
		var found int
	nextItem:
		for i, item := range params.List.Data.Items {
			for _, id := range request.ItemIDs {
				if item.ID == id {
					found++
					continue nextItem
				}
			}
			params.List.Data.Items[i-found] = item
		}
		if found > 0 {
			params.List.Data.Items = params.List.Data.Items[:len(params.List.Data.Items)-found]
		}
		params.ListUpdates = []dal.Update{
			{
				Field: "items",
				Value: params.List.Data.Items,
			},
			{
				Field: "count",
				Value: len(params.List.Data.Items),
			},
		}
		return nil
	})
	if err != nil {
		return err
	}
	return
}
