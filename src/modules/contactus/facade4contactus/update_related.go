package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/strongo/validation"
)

func updateRelatedField(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	objectRef models4linkage.TeamModuleItemRef,
	request dto4contactus.UpdateRelatedRequest,
	params *dal4contactus.ContactWorkerParams, // TODO: needs abstraction
) (err error) {
	relatableAdapted := facade4linkage.NewRelatableAdapter[*models4contactus.ContactDbo](func(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleItemRef) (err error) {
		// Verify contactID belongs to the same team
		teamContactBriefID := recordRef.ItemID
		if _, existingContact := params.TeamModuleEntry.Data.Contacts[teamContactBriefID]; !existingContact {
			if _, err = GetContactByID(ctx, tx, params.Team.ID, recordRef.ItemID); err != nil {
				return fmt.Errorf("failed to get related contact: %w", err)
			}
		}
		return nil
	})
	var itemUpdates, teamModuleUpdates []dal.Update

	for itemID, itemRolesCommand := range request.Related {
		itemRef := models4linkage.NewTeamModuleDocRefFromString(itemID)
		if objectRef == itemRef {
			return validation.NewErrBadRequestFieldValue("itemRef", fmt.Sprintf("objectRef and itemRef are the same: %+v", objectRef))
		}
		if itemUpdates, teamModuleUpdates, err = facade4linkage.SetRelated(ctx, tx, relatableAdapted, params.Contact, objectRef, itemRef, *itemRolesCommand); err != nil {
			return err
		}
		params.ContactUpdates = append(params.ContactUpdates, itemUpdates...)
		params.TeamModuleUpdates = append(params.TeamModuleUpdates, teamModuleUpdates...)
		if err := updateRelatedItem(ctx, tx, itemRef, nil); err != nil {
			return fmt.Errorf("failed to update related record: %w", err)
		}
	}

	return nil
}

func updateRelatedItem(ctx context.Context, tx dal.ReadwriteTransaction, objectRef models4linkage.TeamModuleItemRef, related map[string]*models4linkage.RelationshipRolesCommand) error {
	itemKey := dal4teamus.NewTeamModuleItemKeyFromItemRef(objectRef)
	object := record.NewDataWithID(objectRef.ItemID, itemKey, new(models4linkage.WithRelatedAndIDs))
	if err := tx.Get(ctx, object.Record); err != nil {
		return fmt.Errorf("failed to get object record: %w", err)
	}
	if err := object.Data.Validate(); err != nil {
		return fmt.Errorf("record is not valid after loading from DB: %w", err)
	}
	for itemID /*, itemRolesCommand*/ := range related {
		itemRef := models4linkage.NewTeamModuleDocRefFromString(itemID)
		if objectRef == itemRef {
			return validation.NewErrBadRequestFieldValue("itemRef", fmt.Sprintf("objectRef and itemRef are the same: %+v", objectRef))
		}
		/*
			if itemUpdates, teamModuleUpdates, err := facade4linkage.SetRelated(ctx, tx, relatableAdapted, params.Contact, objectRef, itemRef, *itemRolesCommand); err != nil {
				return err
			}
		*/
	}
	return nil
}
