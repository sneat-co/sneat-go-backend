package facade4listus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/validation"
)

// SetListItemsIsDone marks list item as completed
func SetListItemsIsDone(ctx context.Context, userContext facade.User, request ListItemsSetIsDoneRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	uid := userContext.GetID()
	if uid == "" {
		return validation.NewErrRequestIsMissingRequiredField("userContext.ContactID()")
	}
	db := facade.GetDatabase(ctx)
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		listID := request.ListID
		list := dal4listus.NewSpaceListEntry(request.SpaceID, listID)

		if err := GetListForUpdate(ctx, tx, list); err != nil {
			if dal.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("failed to get a list for update of list items: %w", err)
		}
		//listUpdates := make([]dal.Update, 0, len(request.ItemIDs))
		changed := 0
		for _, item := range list.Data.Items {
			for _, id := range request.ItemIDs {
				logus.Debugf(ctx, "item.InviteID: %s, %s, %s", item.ID, "requestItemID", id)
				if item.ID == id && item.IsDone != request.IsDone {
					item.IsDone = request.IsDone
					changed++
					//listUpdates = append(listUpdates, )
				}
			}
		}
		logus.Debugf(ctx, "changed: %d", changed)
		if changed == 0 {
			return nil
		}
		listUpdates := []dal.Update{
			{
				Field: "items",
				Value: list.Data.Items,
			},
		}
		listKey := list.Record.Key()
		//logus.Debugf(ctx, "Updating list with listKey=%v, item[0]: %+v; updates[0]: %+v",
		//	listKey, list.Data.Items[0], listUpdates[0].Value)
		if err = tx.Update(ctx, listKey, listUpdates); err != nil {
			return fmt.Errorf("failed to update list record: %w", err)
		}
		logus.Debugf(ctx, "Updated list with listKey=%v, field=%s, item[1]: %+v", listKey, listUpdates[0].Field, listUpdates[0].Value)
		return nil
	})
	if err != nil {
		return err
	}
	return
}
