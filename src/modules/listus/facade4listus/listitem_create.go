package facade4listus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// CreateListItems creates list items
func CreateListItems(ctx context.Context, userCtx facade.UserContext, request dto4listus.CreateListItemsRequest) (response dto4listus.CreateListItemResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4spaceus.RunModuleSpaceWorker(ctx, userCtx, request.SpaceID, const4listus.ModuleID, new(dbo4listus.ListusSpaceDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4listus.ListusSpaceDbo]) error {
			return createListItemTxWorker(ctx, request, tx, params)
		})
	return
}

func createListItemTxWorker(ctx context.Context, request dto4listus.CreateListItemsRequest, tx dal.ReadwriteTransaction, params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4listus.ListusSpaceDbo]) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	//if slice.Index(params.Space.Data.UserIDs, uid) < 0 {
	//	// TODO: check if user is a member of the team at RunModuleSpaceWorker() level
	//	return fmt.Errorf("user have no access to this team")
	//}

	listType := request.ListID.ListType()
	var list = dal4listus.NewSpaceListEntry(request.SpaceID, request.ListID)
	if err = tx.Get(ctx, list.Record); err != nil && !dal.IsNotFound(err) {
		return fmt.Errorf("failed to get list record: %w", err)
	}

	if !list.Record.Exists() {

		isOkToAutoCreateList :=
			request.ListID == dbo4listus.NewListKey(dbo4listus.ListTypeToBuy, "groceries") ||
				request.ListID == dbo4listus.NewListKey(dbo4listus.ListTypeToWatch, "movies")

		if !isOkToAutoCreateList {
			err = fmt.Errorf("list not found by ContactID=%s: %w", request.ListID, err)
			return err
		}

		list.Data.SpaceIDs = []string{request.SpaceID}
		list.Data.UserIDs = []string{params.UserID}
		list.Data.Type = listType
		list.Data.Title = string(request.ListID)
		if list.Data.Emoji == "" {
			switch request.ListID.ListType() {
			case dbo4listus.ListTypeToBuy:
				list.Data.Emoji = "🛒"
			case dbo4listus.ListTypeToWatch:
				list.Data.Emoji = "📽️"
			}
		}
	}

	listBrief, isExistingBrief := params.SpaceModuleEntry.Data.Lists[string(request.ListID)]
	if !isExistingBrief {
		params.SpaceModuleEntry.Data.Lists = make(dbo4listus.ListBriefs, 1)
		listBrief = &dbo4listus.ListBrief{
			ListBase: dbo4listus.ListBase{
				Type:  request.ListID.ListType(),
				Title: string(request.ListID),
			},
		}
		if listBrief.Type == dbo4listus.ListTypeToBuy && request.ListID == "groceries" {
			listBrief.Emoji = "🛒"
		}
		params.SpaceModuleEntry.Data.Lists[string(request.ListID)] = listBrief
	}

	for i, item := range request.Items {
		id, err := generateRandomListItemID(list.Data.Items, item.ID)
		if err != nil {
			return fmt.Errorf("failed to generate random id for item #%d: %w", i, err)
		}
		listItem := dbo4listus.ListItemBrief{
			ID:           id,
			ListItemBase: item.ListItemBase,
		}
		if listItem.Emoji == "" {
			listItem.Emoji = deductListItemEmoji(listItem.Title)
		}
		listItem.CreatedAt = params.Started
		listItem.CreatedBy = params.UserID
		list.Data.Items = append(list.Data.Items, &listItem)
	}
	list.Data.Count = len(list.Data.Items)
	listBrief.ItemsCount = len(list.Data.Items)
	if err = list.Data.Validate(); err != nil {
		return fmt.Errorf("list record is not valid: %w", err)
	}
	if list.Record.Exists() {
		if slice.Index(list.Data.UserIDs, params.UserID) < 0 {
			return errors.New("current user does not have access to the list: userID=" + params.UserID)
		}
		if err = tx.Update(ctx, list.Key, []dal.Update{
			{
				Field: "items",
				Value: list.Data.Items,
			},
			{
				Field: "count",
				Value: len(list.Data.Items),
			},
		}); err != nil {
			return fmt.Errorf("failed to update list record: %w", err)
		}
	} else {
		if err = tx.Insert(ctx, list.Record); err != nil {
			return fmt.Errorf("failed to insert list record: %w", err)
		}
	}

	if params.SpaceModuleEntry.Record.Exists() {
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, dal.Update{
			Field: "lists." + string(request.ListID),
			Value: listBrief,
		})
		params.SpaceModuleEntry.Record.MarkAsChanged()
	} else {
		params.SpaceModuleEntry.Data.CreatedAt = params.Started
		params.SpaceModuleEntry.Data.CreatedBy = params.UserID
		if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
			return fmt.Errorf("failed to insert team module entry record: %w", err)
		}
	}

	return nil
}
