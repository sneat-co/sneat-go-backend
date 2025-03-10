package facade4listus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"strings"
)

// CreateListItems creates list items
func CreateListItems(ctx facade.ContextWithUser, request dto4listus.CreateListItemsRequest) (
	response dto4listus.CreateListItemResponse, list dal4listus.ListEntry, err error,
) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4listus.RunListWorker(ctx, ctx.User(), request.ListRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4listus.ListWorkerParams) (err error) {
		response, list, err = createListItemsTxWorker(ctx, tx, request, params)
		if err != nil {
			return fmt.Errorf("failed in createListItemsTxWorker: %w", err)
		}
		return err
	})
	return
}

func createListItemsTxWorker(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4listus.CreateListItemsRequest,
	params *dal4listus.ListWorkerParams,
) (
	response dto4listus.CreateListItemResponse,
	list dal4listus.ListEntry,
	err error,
) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return
	}
	//if slice.Index(params.Space.Data.UserIDs, uid) < 0 {
	//	// TODO: check if user is a member of the team at RunModuleSpaceWorker() level
	//	return fmt.Errorf("user have no access to this team")
	//}

	listType := request.ListID.ListType()
	list = params.List

	if !list.Record.Exists() {

		if !dbo4listus.IsStandardList(request.ListID) {
			err = fmt.Errorf("list not found by listID=%s: %w", request.ListID, err)
			return
		}

		list.Data.SpaceIDs = []coretypes.SpaceID{request.SpaceID}
		list.Data.UserIDs = []string{params.UserID()}
		list.Data.Type = listType
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
		var id string
		if id, err = generateRandomListItemID(list.Data.Items, item.ID); err != nil {
			err = fmt.Errorf("failed to generate random id for item #%d: %w", i, err)
			return
		}
		listItem := &dbo4listus.ListItemBrief{
			ID:           id,
			ListItemBase: item.ListItemBase,
		}
		if listItem.Emoji == "" {
			listItem.Emoji = deductListItemEmoji(listItem.Title)
			if listItem.Emoji != "" && strings.HasPrefix(listItem.Title, listItem.Emoji) {
				listItem.Title = strings.TrimSpace(listItem.Title[len(listItem.Emoji):])
			}
		}
		listItem.CreatedAt = params.Started
		listItem.CreatedBy = params.UserID()
		listItem = list.Data.AddListItem(listItem)
		response.CreatedItems = append(response.CreatedItems, listItem)
	}
	list.Data.Count = len(list.Data.Items)
	listBrief.ItemsCount = len(list.Data.Items)
	if err = list.Data.Validate(); err != nil {
		err = fmt.Errorf("list record is not valid: %w", err)
		return
	}
	if list.Record.Exists() {
		updates := []update.Update{
			update.ByFieldName("items", list.Data.Items),
			update.ByFieldName("count", list.Data.Count),
		}
		userID := params.UserID()
		if !list.Data.HasUserID(userID) {
			if params.Space.Data.HasUserID(params.UserID()) {
				updates = append(updates, list.Data.AddUserID(userID)...)
			} else {
				err = errors.New("current user does not have access to the list: userID=" + params.UserID())
				return
			}
		}
		if err = tx.Update(ctx, list.Key, updates); err != nil {
			err = fmt.Errorf("failed to update list record: %w", err)
			return
		}
	} else {
		if err = tx.Insert(ctx, list.Record); err != nil {
			err = fmt.Errorf("failed to insert list record: %w", err)
			return
		}
	}

	if params.SpaceModuleEntry.Record.Exists() {
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
			update.ByFieldName("lists."+string(request.ListID), listBrief))
		params.SpaceModuleEntry.Record.MarkAsChanged()
	} else {
		params.SpaceModuleEntry.Data.CreatedAt = params.Started
		params.SpaceModuleEntry.Data.CreatedBy = params.UserID()
		if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
			err = fmt.Errorf("failed to insert team module entry record: %w", err)
			return
		}
	}
	return
}
