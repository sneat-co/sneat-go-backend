package facade4listus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/models4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"log"
)

// ReorderListItem reorders list items
func ReorderListItem(ctx context.Context, userContext facade.User, request ReorderListItemsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	uid := userContext.GetID()
	if uid == "" {
		return validation.NewErrRequestIsMissingRequiredField("userContext.ContactID()")
	}
	db := facade.GetDatabase(ctx)
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		listID := models4listus.GetFullListID(request.ListType, request.ListID)
		list := dal4listus.NewTeamListContext(request.TeamID, listID)

		if err = GetListForUpdate(ctx, tx, list); err != nil {
			return fmt.Errorf("failed to get a list for reordering of list items: %w", err)
		}
		//listUpdates := make([]dal.Update, 0, len(request.ItemIDs))
		itemsToMove := make([]*models4listus.ListItemBrief, 0, len(request.ItemIDs))
		otherItems := make([]*models4listus.ListItemBrief, 0, len(list.Dto.Items)-len(request.ItemIDs))

		for _, item := range list.Dto.Items {
			isItemToMove := false
			for _, id := range request.ItemIDs {
				log.Println("item.InviteID", item.ID, "requestItemID", id)
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
		if toIndex >= len(list.Dto.Items) {
			toIndex = len(list.Dto.Items) - 1
		} else if len(otherItems) < toIndex {
			toIndex = len(otherItems) - 1
		}
		items := make([]*models4listus.ListItemBrief, toIndex, len(list.Dto.Items))
		for i := 0; i < toIndex; i++ {
			items[i] = otherItems[i]
		}
		for i := 0; i < len(itemsToMove); i++ {
			items = append(items, itemsToMove[0])
		}
		for i := toIndex; i < len(otherItems); i++ {
			items = append(items, otherItems[i])
		}
		list.Dto.Items = items
		listUpdates := []dal.Update{
			{
				Field: "items",
				Value: list.Dto.Items,
			},
		}
		listKey := list.Record.Key()
		log.Printf("Updating list with listKey=%v, item[1]: %+v; updates[0]: %+v",
			listKey, list.Dto.Items[1], listUpdates[0].Value)
		if err = tx.Update(ctx, listKey, listUpdates); err != nil {
			return fmt.Errorf("failed to update list record: %w", err)
		}
		log.Printf("Updated list with listKey=%v, field=%s, item[1]: %+v", listKey, listUpdates[0].Field, listUpdates[0].Value)
		return nil
	})
	if err != nil {
		return err
	}
	return
}
