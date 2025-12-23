package facade4listus

import (
	"context"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// ReorderListItem reorders list items
func ReorderListItem(ctx facade.ContextWithUser, request dto4listus.ReorderListItemsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	userCtx := ctx.User()
	uid := userCtx.GetUserID()
	if uid == "" {
		return validation.NewErrRequestIsMissingRequiredField("userCtx.ContactID()")
	}
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		listID := request.ListID
		list := dal4listus.NewListEntry(request.SpaceID, listID)

		if err = dal4listus.GetListForUpdate(ctx, tx, list); err != nil {
			return fmt.Errorf("failed to get a list for reordering of list items: %w", err)
		}
		//listUpdates := make([]update.Update, 0, len(request.ItemIDs))
		itemsToMove := make([]*dbo4listus.ListItemBrief, 0, len(request.ItemIDs))
		otherItems := make([]*dbo4listus.ListItemBrief, 0, len(list.Data.Items)-len(request.ItemIDs))

		for _, item := range list.Data.Items {
			isItemToMove := false
			for _, id := range request.ItemIDs {
				if item.ID == id {
					itemsToMove = append(itemsToMove, item)
					isItemToMove = true
					break
				}
			}
			if !isItemToMove {
				otherItems = append(otherItems, item)
			}
		}
		var toIndex = request.ToIndex
		if toIndex >= len(list.Data.Items) {
			toIndex = len(list.Data.Items) - 1
		} else if len(otherItems) < toIndex {
			toIndex = len(otherItems) - 1
		}
		items := make([]*dbo4listus.ListItemBrief, toIndex, len(list.Data.Items))
		for i := 0; i < toIndex; i++ {
			items[i] = otherItems[i]
		}
		for i := 0; i < len(itemsToMove); i++ {
			items = append(items, itemsToMove[0])
		}
		for i := toIndex; i < len(otherItems); i++ {
			items = append(items, otherItems[i])
		}
		list.Data.Items = items
		listUpdates := []update.Update{update.ByFieldName("items", list.Data.Items)}
		listKey := list.Record.Key()
		//logus.Debugf("Updating list with listKey=%v, item[1]: %+v; updates[0]: %+v",
		//	listKey, list.Data.Items[1], listUpdates[0].Value)
		if err = tx.Update(ctx, listKey, listUpdates); err != nil {
			return fmt.Errorf("failed to update list record: %w", err)
		}
		//logus.Debugf("Updated list with listKey=%v, field=%s, item[1]: %+v", listKey, listUpdates[0].Field, listUpdates[0].Value)
		return nil
	})
	if err != nil {
		return err
	}
	return
}
