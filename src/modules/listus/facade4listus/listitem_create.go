package facade4listus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/models4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// CreateListItems creates list items
func CreateListItems(ctx context.Context, userContext facade.User, request CreateListItemsRequest) (response CreateListItemResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	uid := userContext.GetID()
	err = dal4teamus.RunModuleTeamWorker(ctx, userContext, request.TeamRequest, const4listus.ModuleID, new(models4listus.ListusTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4listus.ListusTeamDto]) error {
			return createListItemTxWorker(ctx, request, uid, tx, params)
		})
	return
}

func createListItemTxWorker(ctx context.Context, request CreateListItemsRequest, uid string, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4listus.ListusTeamDto]) (err error) {
	if err = params.GetRecords(ctx, tx, uid); err != nil {
		return err
	}
	//if slice.Index(params.Team.Data.UserIDs, uid) < 0 {
	//	// TODO: check if user is a member of the team at RunModuleTeamWorker() level
	//	return fmt.Errorf("user have no access to this team")
	//}

	listID := models4listus.GetFullListID(request.ListType, request.ListID)
	listKey := dal4listus.NewTeamListKey(request.TeamID, listID)
	var listDto models4listus.ListDto
	var listRecord = dal.NewRecordWithData(listKey, &listDto)
	if err = tx.Get(ctx, listRecord); err != nil && !dal.IsNotFound(err) {
		return fmt.Errorf("failed to get list record: %w", err)
	}

	if !listRecord.Exists() {
		isOkToAutoCreateList :=
			request.ListType == models4listus.ListTypeToBuy && request.ListID == "groceries" ||
				request.ListType == models4listus.ListTypeToWatch && request.ListID == "movies"

		if !isOkToAutoCreateList {
			err = fmt.Errorf("list not found by ID=%s: %w", listID, err)
			return err
		}

		listDto.TeamIDs = []string{request.TeamID}
		listDto.UserIDs = []string{uid}
		listDto.Type = request.ListType
		listDto.Title = request.ListID
		if request.ListType == "to-buy" && request.ListID == "groceries" {
			listDto.Emoji = "🛒"
		}
	}

	listBrief, isExistingBrief := params.TeamModuleEntry.Data.Lists[listID]
	if !isExistingBrief {
		params.TeamModuleEntry.Data.Lists = make(map[string]*models4listus.ListBrief, 1)
		listBrief = &models4listus.ListBrief{
			ListBase: models4listus.ListBase{
				Type:  request.ListType,
				Title: request.ListID,
			},
		}
		if listBrief.Type == models4listus.ListTypeToBuy && request.ListID == "groceries" {
			listBrief.Emoji = "🛒"
		}
		params.TeamModuleEntry.Data.Lists[listID] = listBrief
	}

	for i, item := range request.Items {
		id, err := generateRandomListItemID(listDto.Items, item.ID)
		if err != nil {
			return fmt.Errorf("failed to generate random id for item #%v: %w", i, err)
		}
		listItem := models4listus.ListItemBrief{
			ID:           id,
			ListItemBase: item.ListItemBase,
		}
		listItem.CreatedAt = params.Started
		listItem.CreatedBy = uid
		listDto.Items = append(listDto.Items, &listItem)
	}
	listDto.Count = len(listDto.Items)
	listBrief.ItemsCount = len(listDto.Items)
	if err := listDto.Validate(); err != nil {
		return fmt.Errorf("list record is not valid: %w", err)
	}
	if listRecord.Exists() {
		if slice.Index(listDto.UserIDs, uid) < 0 {
			return errors.New("current user does not have access to the list: userID=" + uid)
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

	if params.TeamModuleEntry.Record.Exists() {
		params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{
			Field: "lists." + listID,
			Value: listBrief,
		})
	} else {
		params.TeamModuleEntry.Data.CreatedAt = params.Started
		params.TeamModuleEntry.Data.CreatedBy = uid
		if err = tx.Insert(ctx, params.TeamModuleEntry.Record); err != nil {
			return fmt.Errorf("failed to insert team module entry record: %w", err)
		}
	}

	return nil
}
