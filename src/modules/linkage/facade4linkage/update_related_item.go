package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/strongo/validation"
)

func updateRelatedItem(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	objectRef dbo4linkage.SpaceModuleItemRef,
	relateds []*dbo4linkage.RelationshipItemRolesCommand,
) (recordsUpdates []dal4teamus.RecordUpdates, err error) {
	itemKey := dal4teamus.NewSpaceModuleItemKeyFromItemRef(objectRef)
	itemDbo := new(dbo4linkage.WithRelatedAndIDsAndUserID)
	itemDbo.WithRelatedAndIDs = new(dbo4linkage.WithRelatedAndIDs)
	object := record.NewDataWithID(objectRef.ItemID, itemKey, itemDbo)
	if err = tx.Get(ctx, object.Record); err != nil {
		return recordsUpdates, fmt.Errorf("failed to get object record: %w", err)
	}
	if err = object.Data.Validate(); err != nil {
		return recordsUpdates, fmt.Errorf("record is not valid after loading from DB: %w", err)
	}
	for _, related := range relateds {
		if objectRef == related.ItemRef {
			return recordsUpdates, validation.NewErrBadRequestFieldValue("itemRef", fmt.Sprintf("objectRef and itemRef are the same: %+v", objectRef))
		}
		/*
			if itemUpdates, teamModuleUpdates, err := facade4linkage.SetRelated(ctx, tx, relatableAdapted, params.ContactEntry, objectRef, itemRef, *itemRolesCommand); err != nil {
				return err
			}
		*/
	}
	if object.Data.UserID != "" {
		if userUpdates, err := updateUserRelated(ctx, tx,
			object.Data.UserID, objectRef,
			record.NewDataWithID(object.ID, object.Key, &object.Data.WithRelated),
		); err != nil {
			key := dal4userus.NewUserModuleKey(object.Data.UserID, objectRef.Module)
			return recordsUpdates, fmt.Errorf("failed to update related field in %s: %w", key.String(), err)
		} else {
			recordsUpdates = append(recordsUpdates, userUpdates)
		}
	}
	return recordsUpdates, nil
}
