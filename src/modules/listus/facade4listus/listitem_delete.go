package facade4listus

import (
	"slices"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// DeleteListItems deletes list items
func DeleteListItems(ctx facade.ContextWithUser, request dto4listus.ListItemIDsRequest) (deletedItems []*dbo4listus.ListItemBrief, list dal4listus.ListEntry, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4listus.RunListWorker(ctx, request.ListRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4listus.ListWorkerParams) (err error) {
			list = params.List
			isInRecentItems := func(item *dbo4listus.ListItemBrief) bool {
				for _, recentItem := range list.Data.RecentItems {
					if recentItem.ID == item.ID || recentItem.Title == item.Title && recentItem.Emoji == item.Emoji {
						return true
					}
				}
				return false
			}
			removeAll := len(request.ItemIDs) == 1 && request.ItemIDs[0] == "*"
			var recentItems []*dbo4listus.ListItemBrief
			items, removedCount := slice.RemoveInPlace(params.List.Data.Items, func(item *dbo4listus.ListItemBrief) (remove bool) {
				if remove = removeAll || slices.Contains(request.ItemIDs, item.ID); remove {
					deletedItems = append(deletedItems, item)
					if !isInRecentItems(item) {
						recentItems = append(recentItems, item)
					}
				}
				return
			})
			if removedCount > 0 {
				params.List.Record.MarkAsChanged()
				params.List.Data.Items = items
				params.List.Data.Count = len(items)
				params.ListUpdates = append(params.ListUpdates,
					update.ByFieldName("items", params.List.Data.Items),
					update.ByFieldName("count", len(params.List.Data.Items)),
				)
				if len(recentItems) > 0 {
					slices.Reverse(recentItems)
					prevRecentItems := list.Data.RecentItems
					list.Data.RecentItems = recentItems
					for _, prevRecentItem := range prevRecentItems {
						if len(list.Data.RecentItems) >= 100 {
							break
						}
						list.Data.RecentItems = append(list.Data.RecentItems, prevRecentItem)
					}

					// This should be after we set the params.ListUpdates
					params.ListUpdates = append(params.ListUpdates, update.ByFieldName("recentItems", list.Data.RecentItems))
				}
			}
			return
		})
	return
}
