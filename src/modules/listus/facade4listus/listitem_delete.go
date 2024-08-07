package facade4listus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// DeleteListItems deletes list items
func DeleteListItems(ctx context.Context, userContext facade.User, request ListItemIDsRequest) (err error) {
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
			return err
		}
		var found int
	nextItem:
		for i, item := range list.Data.Items {
			for _, id := range request.ItemIDs {
				if item.ID == id {
					found++
					continue nextItem
				}
			}
			list.Data.Items[i-found] = item
		}
		if found > 0 {
			list.Data.Items = list.Data.Items[:len(list.Data.Items)-found]
		}
		listUpdates := []dal.Update{
			{
				Field: "items",
				Value: list.Data.Items,
			},
			{
				Field: "count",
				Value: len(list.Data.Items),
			},
		}
		if err = tx.Update(ctx, list.Record.Key(), listUpdates); err != nil {
			return fmt.Errorf("failed to update list record: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return
}
