package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/strongo/validation"
)

func updateRelatedItem(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	objectRef models4linkage.TeamModuleItemRef,
	related map[string]*models4linkage.RelationshipRolesCommand,
) (recordsUpdates []dal4teamus.RecordUpdates, err error) {
	itemKey := dal4teamus.NewTeamModuleItemKeyFromItemRef(objectRef)
	object := record.NewDataWithID(objectRef.ItemID, itemKey, new(models4linkage.WithRelatedAndIDsAndUserID))
	if err := tx.Get(ctx, object.Record); err != nil {
		return recordsUpdates, fmt.Errorf("failed to get object record: %w", err)
	}
	if err := object.Data.Validate(); err != nil {
		return recordsUpdates, fmt.Errorf("record is not valid after loading from DB: %w", err)
	}
	for itemID /*, itemRolesCommand*/ := range related {
		itemRef := models4linkage.NewTeamModuleItemRefFromString(itemID)
		if objectRef == itemRef {
			return recordsUpdates, validation.NewErrBadRequestFieldValue("itemRef", fmt.Sprintf("objectRef and itemRef are the same: %+v", objectRef))
		}
		/*
			if itemUpdates, teamModuleUpdates, err := facade4linkage.SetRelated(ctx, tx, relatableAdapted, params.Contact, objectRef, itemRef, *itemRolesCommand); err != nil {
				return err
			}
		*/
	}
	if object.Data.UserID != "" {
		if userUpdates, err := updateUserRelated(ctx, tx,
			object.Data.UserID, objectRef,
			record.NewDataWithID(object.ID, object.Key, &object.Data.WithRelated),
		); err != nil {
			key := dal4userus.NewUserModuleKey(object.Data.UserID, objectRef.ModuleID)
			return recordsUpdates, fmt.Errorf("failed to update related field in %s: %w", key.String(), err)
		} else {
			recordsUpdates = append(recordsUpdates, userUpdates)
		}
	}
	return recordsUpdates, nil
}
