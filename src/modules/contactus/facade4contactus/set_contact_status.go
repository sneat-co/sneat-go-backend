package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// SetContactsStatus sets contacts status
func SetContactsStatus(ctx context.Context, user facade.User, request dto4contactus.SetContactsStatusRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	err = dal4contactus.RunContactusTeamWorker(ctx, user, request.TeamRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusTeamWorkerParams) (err error) {
			return setContactsStatusTxWorker(ctx, tx, params, request.ContactIDs, request.Status)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to set contact status: %w", err)
	}
	return nil
}

func setContactsStatusTxWorker(
	ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusTeamWorkerParams,
	contactIDs []string, status string,
) (err error) {
	for _, contactID := range contactIDs {
		if err := setContactStatusTxWorker(ctx, tx, params, contactID, status); err != nil {
			return fmt.Errorf("failed to set status for contact id=[%s]: %w", contactID, err)
		}
	}
	return nil
}

func setContactStatusTxWorker(
	ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusTeamWorkerParams,
	contactID string, status string,
) (err error) {
	contact := dal4contactus.NewContactEntry(params.Team.ID, contactID)
	if err = tx.Get(ctx, contact.Record); err != nil {
		return fmt.Errorf("failed to get contact record: %w", err)
	}

	var relatedContacts []dal4contactus.ContactEntry

	relatedContacts, err = GetRelatedContacts(ctx, tx, params.Team.ID, "child", 0, -1, []dal4contactus.ContactEntry{contact})
	if err != nil {
		return fmt.Errorf("failed to get descendant contacts: %w", err)
	}
	contactsToUpdate := append(relatedContacts, contact)
	contactKeys := make([]*dal.Key, 0, len(relatedContacts)+1)
	for _, contactToUpdate := range contactsToUpdate {
		if contactToUpdate.Data.Status != status {
			contactToUpdate.Data.Status = status
			contactKeys = append(contactKeys, contactToUpdate.Key)
			if err := contact.Data.Validate(); err != nil {
				return err
			}
		}
	}
	if len(contactKeys) > 0 {
		if err := tx.UpdateMulti(ctx, contactKeys, []dal.Update{
			{Field: "status", Value: status},
		}); err != nil {
			return fmt.Errorf("failed to update contact records to set status to %v: %w", status, err)
		}
	}
	if status == dbmodels.StatusArchived || status == dbmodels.StatusDeleted {
		contactIDs := make([]string, 0, len(contactsToUpdate))
		for _, contactToUpdate := range contactsToUpdate {
			contactIDs = append(contactIDs, contactToUpdate.ID)
		}
		for _, contactID := range contactIDs {
			params.TeamModuleUpdates = append(params.TeamModuleUpdates,
				params.TeamModuleEntry.Data.RemoveContact(contactID))
		}
		if err := params.Team.Data.Validate(); err != nil {
			return err
		}
		//params.TeamUpdates = append(params.TeamUpdates, updateTeamDtoWithNumberOfContact(len(params.TeamModuleEntry.Data.Contacts)))
	}
	if status == "active" {
		params.TeamModuleEntry.Data.AddContact(contact.ID, &contact.Data.ContactBrief)
	}
	if params.TeamModuleEntry.Record.Exists() {
		if len(params.TeamModuleEntry.Data.Contacts) == 0 {
			if err := tx.Delete(ctx, params.TeamModuleEntry.Key); err != nil {
				return fmt.Errorf("failed to delete team contacts brief record: %w", err)
			}
		} else {
			if err := tx.Update(ctx, params.TeamModuleEntry.Key, []dal.Update{
				{
					Field: const4contactus.ContactsField,
					Value: params.TeamModuleEntry.Data.Contacts,
				},
			}); err != nil {
				return fmt.Errorf("failed to put team contacts brief: %w", err)
			}
		}
	} else if len(params.TeamModuleEntry.Data.Contacts) > 0 {
		if err := tx.Insert(ctx, params.TeamModuleEntry.Record); err != nil {
			return fmt.Errorf("failed to insert team contacts brief record: %w", err)
		}
	}
	return nil
}
