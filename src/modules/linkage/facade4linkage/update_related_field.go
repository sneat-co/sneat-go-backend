package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/strongo/validation"
)

func UpdateRelatedField(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	objectRef dbo4linkage.SpaceModuleItemRef,
	request dto4linkage.UpdateRelatedFieldRequest,
	item *dbo4linkage.WithRelatedAndIDsAndUserID,
	addUpdates func(updates []dal.Update),
) (recordsUpdates []dal4teamus.RecordUpdates, err error) {
	var setRelatedResult SetRelatedResult

	for itemID, itemRolesCommand := range request.Related {
		itemRef := dbo4linkage.NewSpaceModuleItemRefFromString(itemID)
		if objectRef == itemRef {
			return recordsUpdates, validation.NewErrBadRequestFieldValue("itemRef", fmt.Sprintf("objectRef and itemRef are the same: %+v", objectRef))
		}
		if setRelatedResult, err = SetRelated(ctx, tx,
			item, objectRef, itemRef, *itemRolesCommand); err != nil {
			return recordsUpdates, err
		}

		addUpdates(setRelatedResult.ItemUpdates)
		//params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, setRelatedResult.SpaceModuleUpdates...)

		if recordsUpdates, err = updateRelatedItem(ctx, tx, itemRef, nil); err != nil {
			return recordsUpdates, fmt.Errorf("failed to update related record: %w", err)
		}
	}

	return recordsUpdates, nil
}
