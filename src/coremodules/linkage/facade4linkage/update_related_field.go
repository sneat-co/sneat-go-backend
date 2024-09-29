package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	dbo4linkage2 "github.com/sneat-co/sneat-go-backend/src/coremodules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/linkage/dto4linkage"
	"github.com/strongo/validation"
)

func UpdateRelatedField(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	objectRef dbo4linkage2.SpaceModuleItemRef,
	request dto4linkage.UpdateRelatedFieldRequest,
	item *dbo4linkage2.WithRelatedAndIDsAndUserID,
	addUpdates func(updates []dal.Update),
) (recordsUpdates []record.Updates, err error) {
	var setRelatedResult SetRelatedResult

	for i, itemRolesCommand := range request.Related {
		itemRef := itemRolesCommand.ItemRef
		if objectRef == itemRef {
			return recordsUpdates, validation.NewErrBadRequestFieldValue("itemRef", fmt.Sprintf("objectRef and itemRef are the same: %+v", objectRef))
		}
		if setRelatedResult, err = SetRelated(ctx, tx, item, objectRef, itemRef, itemRolesCommand); err != nil {
			return recordsUpdates, err
		}

		addUpdates(setRelatedResult.ItemUpdates)
		//params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, setRelatedResult.SpaceModuleUpdates...)

		if recordsUpdates, err = updateRelatedItem(ctx, tx, itemRef, nil); err != nil {
			return recordsUpdates, fmt.Errorf("failed to update related record for command [%d=%s]: %w", i, itemRef.ID(), err)
		}
	}

	return recordsUpdates, nil
}
