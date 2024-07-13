package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateItemRelationships(ctx context.Context, userCtx facade.User, request dto4linkage.UpdateItemRequest) (item record.DataWithID[string, *dbo4linkage.WithRelatedAndIDsAndUserID], err error) {
	if err = dal4teamus.RunSpaceWorker(ctx, userCtx, request.Space, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.SpaceWorkerParams) (err error) {
		item, err = txUpdateItemRelationships(ctx, tx, params, request)
		return err
	}); err != nil {
		return item, err
	}
	if err = UpdateRelatedItemsWithLatestRelationships(ctx, userCtx, request, *item.Data.WithRelatedAndIDs); err != nil {
		return item, err
	}
	return item, err
}

func txUpdateItemRelationships(
	ctx context.Context, tx dal.ReadwriteTransaction,
	params *dal4teamus.SpaceWorkerParams,
	request dto4linkage.UpdateItemRequest,
) (item record.DataWithID[string, *dbo4linkage.WithRelatedAndIDsAndUserID], err error) {
	key := dal4teamus.NewSpaceModuleItemKey(request.Space, request.Module, request.Collection, request.ItemID)
	item = record.NewDataWithID[string, *dbo4linkage.WithRelatedAndIDsAndUserID](request.ItemID, key, new(dbo4linkage.WithRelatedAndIDsAndUserID))
	if err = tx.Get(ctx, item.Record); err != nil {
		return item, err
	}
	var itemUpdates []dal.Update
	params.RecordUpdates, err = UpdateRelatedField(ctx, tx,
		request.SpaceModuleItemRef, request.UpdateRelatedFieldRequest, item.Data,
		func(updates []dal.Update) {
			itemUpdates = append(itemUpdates, updates...)
		})
	if err != nil {
		return item, err
	}
	if len(itemUpdates) > 0 {
		if err = tx.Update(ctx, item.Key, itemUpdates); err != nil {
			return item, err
		}
	}
	return item, nil
}
