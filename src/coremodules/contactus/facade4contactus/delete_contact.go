package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// DeleteContact deletes team contact
func DeleteContact(ctx context.Context, userCtx facade.UserContext, request dto4contactus.ContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	return dal4contactus.RunContactusSpaceWorker(ctx, userCtx, request.SpaceRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) (err error) {
			return deleteContactTxWorker(ctx, tx, params, request.ContactID)
		},
	)
}

func deleteContactTxWorker(
	ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams,
	contactID string,
) (err error) {
	if contactID == params.Space.ID {
		return validation.NewErrBadRequestFieldValue("contactID", "cannot delete contact that represents team/company itself")
	}
	contact := dal4contactus.NewContactEntry(params.Space.ID, contactID)
	if err = params.GetRecords(ctx, tx, params.Space.Record); err != nil {
		return err
	}

	var subContacts []dal4contactus.ContactEntry
	subContacts, err = GetRelatedContacts(ctx, tx, params.Space.ID, RelatedAsChild, 0, -1, []dal4contactus.ContactEntry{contact})
	if err != nil {
		return fmt.Errorf("failed to get related contacts: %w", err)
	}

	params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
		params.SpaceModuleEntry.Data.RemoveContact(contactID))

	if err := params.Space.Data.Validate(); err != nil {
		return err
	}

	//params.SpaceUpdates = append(params.SpaceUpdates, updateTeamDtoWithNumberOfContact(len(params.SpaceModuleEntry.Data.Contacts)))

	contactKeysToDelete := make([]*dal.Key, 0, len(subContacts)+1)
	contactKeysToDelete = append(contactKeysToDelete, contact.Key)
	for _, subContact := range subContacts {
		subContact.Data.Status = dbmodels.StatusDeleted
		contactKeysToDelete = append(contactKeysToDelete, subContact.Key)
	}
	contactsUpdates := []dal.Update{
		{
			Field: "status",
			Value: dbmodels.StatusDeleted,
		},
	}
	if err = tx.UpdateMulti(ctx, contactKeysToDelete, contactsUpdates); err != nil {
		return fmt.Errorf("failed to set contacts status to %v: %w", contactsUpdates[0].Value, err)
	}
	//if err = tx.DeleteMulti(ctx, contactKeysToDelete); err != nil {
	//	return fmt.Errorf("failed to delete contacts: %w", err)
	//}
	return nil
}
