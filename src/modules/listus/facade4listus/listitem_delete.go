package facade4listus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
	"slices"
)

// DeleteListItems deletes list items
func DeleteListItems(ctx context.Context, userCtx facade.UserContext, request dto4listus.ListItemIDsRequest) (deletedItems []*dbo4listus.ListItemBrief, list dal4listus.ListEntry, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4listus.RunListWorker(ctx, userCtx, request.ListRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4listus.ListWorkerParams) (err error) {
		list = params.List
		items, removedCount := slice.RemoveInPlace(params.List.Data.Items, func(item *dbo4listus.ListItemBrief) bool {
			if slices.Contains(request.ItemIDs, item.ID) {
				deletedItems = append(deletedItems, item)
				return true
			}
			return false
		})
		if removedCount > 0 {
			params.List.Data.Items = items
			params.List.Data.Count = len(items)
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
			params.List.Record.MarkAsChanged()
		}
		return
	})
	return
}
