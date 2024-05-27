package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func updateRelatedField(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	objectRef models4linkage.TeamModuleItemRef,
	request dto4linkage.UpdateRelatedFieldRequest,
	item *WithRelatedAndIDsAndUserID,
	addUpdates func(updates []dal.Update),
) (recordsUpdates []dal4teamus.RecordUpdates, err error) {
	var setRelatedResult facade4linkage.SetRelatedResult

	for itemID, itemRolesCommand := range request.Related {
		itemRef := models4linkage.NewTeamModuleDocRefFromString(itemID)
		if objectRef == itemRef {
			return recordsUpdates, validation.NewErrBadRequestFieldValue("itemRef", fmt.Sprintf("objectRef and itemRef are the same: %+v", objectRef))
		}
		if setRelatedResult, err = facade4linkage.SetRelated(ctx, tx,
			item, objectRef, itemRef, *itemRolesCommand); err != nil {
			return recordsUpdates, err
		}

		addUpdates(setRelatedResult.ItemUpdates)
		//params.TeamModuleUpdates = append(params.TeamModuleUpdates, setRelatedResult.TeamModuleUpdates...)

		if recordsUpdates, err = updateRelatedItem(ctx, tx, itemRef, nil); err != nil {
			return recordsUpdates, fmt.Errorf("failed to update related record: %w", err)
		}
	}

	return recordsUpdates, nil
}

type WithRelatedAndIDsAndUserID struct {
	dbmodels.WithUserID
	*models4linkage.WithRelatedAndIDs
}

func (v *WithRelatedAndIDsAndUserID) Validate() error {
	if err := v.WithUserID.Validate(); err != nil {
		return err
	}
	if err := v.WithRelatedAndIDs.Validate(); err != nil {
		return err
	}
	return nil
}

func updateRelatedItem(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	objectRef models4linkage.TeamModuleItemRef,
	related map[string]*models4linkage.RelationshipRolesCommand,
) (recordsUpdates []dal4teamus.RecordUpdates, err error) {
	itemKey := dal4teamus.NewTeamModuleItemKeyFromItemRef(objectRef)
	object := record.NewDataWithID(objectRef.ItemID, itemKey, new(WithRelatedAndIDsAndUserID))
	if err := tx.Get(ctx, object.Record); err != nil {
		return recordsUpdates, fmt.Errorf("failed to get object record: %w", err)
	}
	if err := object.Data.Validate(); err != nil {
		return recordsUpdates, fmt.Errorf("record is not valid after loading from DB: %w", err)
	}
	for itemID /*, itemRolesCommand*/ := range related {
		itemRef := models4linkage.NewTeamModuleDocRefFromString(itemID)
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

func updateUserRelated(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userID string,
	objectRef models4linkage.TeamModuleItemRef,
	item record.DataWithID[string, *models4linkage.WithRelated],
) (userUpdates dal4teamus.RecordUpdates, err error) {

	user := models4userus.NewUser(userID)
	if err = tx.Get(ctx, user.Record); err != nil {
		return
	}

	return
}
