package facade4linkage

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateRelatedItemsWithLatestRelationships(
	ctx context.Context,
	userCtx facade.User,
	request dto4linkage.UpdateItemRequest,
	itemData dbo4linkage.WithRelatedAndIDs,
) (err error) {
	var updateErrors []error
	for i, related := range request.Related {
		err = updateItemWithLatestRelationshipsFromRelatedItem(ctx, userCtx, related.ItemRef, request.SpaceModuleItemRef, itemData.Related)
		if err != nil {
			updateErrors = append(updateErrors, fmt.Errorf("failed to update related item (%d=%s): %w", i, related.ItemRef.ID(), err))
		}
	}
	if len(updateErrors) > 0 {
		return fmt.Errorf("failed to update %d related items: %w", len(updateErrors), errors.Join(updateErrors...))
	}
	return nil
}

func updateItemWithLatestRelationshipsFromRelatedItem(
	ctx context.Context,
	_ facade.User,
	itemRef dbo4linkage.SpaceModuleItemRef,
	relatedItemRef dbo4linkage.SpaceModuleItemRef,
	relatedByModuleOfRelatedItem dbo4linkage.RelatedByModuleID,
) (err error) {
	itemRelationshipInRelatedItem := dbo4linkage.GetRelatedItemByRef(relatedByModuleOfRelatedItem, itemRef, false)
	if itemRelationshipInRelatedItem == nil || len(itemRelationshipInRelatedItem.RolesToItem) == 0 {
		return nil
	}

	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return err
	}

	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {

		key := dal4spaceus.NewSpaceModuleItemKey(itemRef.Space, itemRef.Module, itemRef.Collection, itemRef.ItemID)
		item := record.NewDataWithID(itemRef.ItemID, key, new(dbo4linkage.WithRelatedAndIDsAndUserID))
		if err = tx.Get(ctx, item.Record); err != nil {
			return err
		}
		relatedItem := dbo4linkage.GetRelatedItemByRef(item.Data.Related, relatedItemRef, true)

		// We do not override existing roles in related item, so we do not lose a role in case of race condition
		for roleID, role := range itemRelationshipInRelatedItem.RolesToItem {
			relatedItem.RolesOfItem[roleID] = role
		}
		if item.Data.UserID != "" {
			users := make(map[string]dbo4userus.User)
			if err = updateUserWithRelatedTx(ctx, tx, item.Data.UserID, users, itemRef, *relatedItem); err != nil {
				return err
			}

		}
		return nil
	})
}
