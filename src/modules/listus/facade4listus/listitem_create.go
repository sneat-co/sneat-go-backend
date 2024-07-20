package facade4listus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// CreateListItems creates list items
func CreateListItems(ctx context.Context, userContext facade.User, request CreateListItemsRequest) (response CreateListItemResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4spaceus.RunModuleSpaceWorker(ctx, userContext, request.SpaceRequest, const4listus.ModuleID, new(dbo4listus.ListusSpaceDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4listus.ListusSpaceDbo]) error {
			return createListItemTxWorker(ctx, request, tx, params)
		})
	return
}

func createListItemTxWorker(ctx context.Context, request CreateListItemsRequest, tx dal.ReadwriteTransaction, params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4listus.ListusSpaceDbo]) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	//if slice.Index(params.Space.Data.UserIDs, uid) < 0 {
	//	// TODO: check if user is a member of the team at RunModuleSpaceWorker() level
	//	return fmt.Errorf("user have no access to this team")
	//}

	listType := request.ListType()
	listID := request.ListID
	listKey := dal4listus.NewSpaceListKey(request.SpaceID, listID)
	var listDto dbo4listus.ListDbo
	var listRecord = dal.NewRecordWithData(listKey, &listDto)
	if err = tx.Get(ctx, listRecord); err != nil && !dal.IsNotFound(err) {
		return fmt.Errorf("failed to get list record: %w", err)
	}

	if !listRecord.Exists() {

		isOkToAutoCreateList :=
			request.ListID == dbo4listus.GetFullListID(dbo4listus.ListTypeToBuy, "groceries") ||
				request.ListID == dbo4listus.GetFullListID(dbo4listus.ListTypeToWatch, "movies")

		if !isOkToAutoCreateList {
			err = fmt.Errorf("list not found by ID=%s: %w", listID, err)
			return err
		}

		listDto.SpaceIDs = []string{request.SpaceID}
		listDto.UserIDs = []string{params.UserID}
		listDto.Type = listType
		listDto.Title = request.ListID
		if listDto.Emoji == "" {
			switch request.ListType() {
			case dbo4listus.ListTypeToBuy:
				listDto.Emoji = "🛒"
			case dbo4listus.ListTypeToWatch:
				listDto.Emoji = "📽️"
			}
		}
	}

	listBrief, isExistingBrief := params.SpaceModuleEntry.Data.Lists[listID]
	if !isExistingBrief {
		params.SpaceModuleEntry.Data.Lists = make(map[string]*dbo4listus.ListBrief, 1)
		listBrief = &dbo4listus.ListBrief{
			ListBase: dbo4listus.ListBase{
				Type:  request.ListType(),
				Title: request.ListID,
			},
		}
		if listBrief.Type == dbo4listus.ListTypeToBuy && request.ListID == "groceries" {
			listBrief.Emoji = "🛒"
		}
		params.SpaceModuleEntry.Data.Lists[listID] = listBrief
	}

	for i, item := range request.Items {
		id, err := generateRandomListItemID(listDto.Items, item.ID)
		if err != nil {
			return fmt.Errorf("failed to generate random id for item #%d: %w", i, err)
		}
		listItem := dbo4listus.ListItemBrief{
			ID:           id,
			ListItemBase: item.ListItemBase,
		}
		listItem.CreatedAt = params.Started
		listItem.CreatedBy = params.UserID
		listDto.Items = append(listDto.Items, &listItem)
	}
	listDto.Count = len(listDto.Items)
	listBrief.ItemsCount = len(listDto.Items)
	if err := listDto.Validate(); err != nil {
		return fmt.Errorf("list record is not valid: %w", err)
	}
	if listRecord.Exists() {
		if slice.Index(listDto.UserIDs, params.UserID) < 0 {
			return errors.New("current user does not have access to the list: userID=" + params.UserID)
		}
		if err := tx.Update(ctx, listKey, []dal.Update{
			{
				Field: "items",
				Value: listDto.Items,
			},
			{
				Field: "count",
				Value: len(listDto.Items),
			},
		}); err != nil {
			return fmt.Errorf("failed to update list record: %w", err)
		}
	} else {
		if err := tx.Insert(ctx, listRecord); err != nil {
			return fmt.Errorf("failed to insert list record: %w", err)
		}
	}

	if params.SpaceModuleEntry.Record.Exists() {
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, dal.Update{
			Field: "lists." + listID,
			Value: listBrief,
		})
	} else {
		params.SpaceModuleEntry.Data.CreatedAt = params.Started
		params.SpaceModuleEntry.Data.CreatedBy = params.UserID
		if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
			return fmt.Errorf("failed to insert team module entry record: %w", err)
		}
	}

	return nil
}
