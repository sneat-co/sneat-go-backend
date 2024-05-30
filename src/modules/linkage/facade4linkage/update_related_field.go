package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/strongo/validation"
)

func UpdateRelatedField(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	objectRef models4linkage.TeamModuleItemRef,
	request dto4linkage.UpdateRelatedFieldRequest,
	item *models4linkage.WithRelatedAndIDsAndUserID,
	addUpdates func(updates []dal.Update),
) (recordsUpdates []dal4teamus.RecordUpdates, err error) {
	var setRelatedResult SetRelatedResult

	for itemID, itemRolesCommand := range request.Related {
		itemRef := models4linkage.NewTeamModuleItemRefFromString(itemID)
		if objectRef == itemRef {
			return recordsUpdates, validation.NewErrBadRequestFieldValue("itemRef", fmt.Sprintf("objectRef and itemRef are the same: %+v", objectRef))
		}
		if setRelatedResult, err = SetRelated(ctx, tx,
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
