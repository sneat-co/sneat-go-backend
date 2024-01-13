package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// UpdateContact sets contact fields
func UpdateContact(
	ctx context.Context,
	user facade.User,
	request dto4contactus.UpdateContactRequest,
) (err error) {
	return RunContactWorker(ctx, user, request.ContactRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactWorkerParams) (err error) {
		return UpdateContactTx(ctx, tx, request, params)
	})
}

// UpdateContactTx sets contact fields
func UpdateContactTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4contactus.UpdateContactRequest,
	params *ContactWorkerParams,
) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	return updateContactTxWorker(ctx, tx, request, params)
}

func updateContactTxWorker(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4contactus.UpdateContactRequest,
	params *ContactWorkerParams,
) (err error) {
	contact := dal4contactus.NewContactEntry(params.Team.ID, request.ContactID)
	//contactData := contact.Data
	//contactData.Validate()

	if err = params.GetRecords(ctx, tx, params.UserID, contact.Record); err != nil {
		return err
	}

	if err := contact.Data.Validate(); err != nil {
		return fmt.Errorf("contact DTO is not valid after loading from DB: %w", err)
	}

	contactBrief := params.TeamModuleEntry.Data.Contacts[request.ContactID]

	var updatedContactFields []string
	var contactUpdates []dal.Update

	if request.Address != nil {
		if *request.Address != *contact.Data.Address {
			updatedContactFields = append(updatedContactFields, "address")
			contact.Data.Address = request.Address
			contactUpdates = append(contactUpdates, dal.Update{Field: "address", Value: request.Address})
		}
	}

	if request.VatNumber != nil {
		if vat := *request.VatNumber; vat != contact.Data.VATNumber {
			updatedContactFields = append(updatedContactFields, "vatNumber")
			contact.Data.VATNumber = vat
			contactUpdates = append(contactUpdates, dal.Update{Field: "vatNumber", Value: vat})
		}
	}

	if request.AgeGroup != "" {
		if request.AgeGroup != contact.Data.AgeGroup {
			updatedContactFields = append(updatedContactFields, "ageGroup")
			contact.Data.AgeGroup = request.AgeGroup
			contactUpdates = append(contactUpdates, dal.Update{Field: "ageGroup", Value: contact.Data.AgeGroup})
		}
		if contactBrief != nil && contactBrief.AgeGroup != request.AgeGroup {
			params.TeamModuleUpdates = append(params.TeamModuleUpdates,
				dal.Update{
					Field: fmt.Sprintf("contacts.%s.ageGroup", request.ContactID),
					Value: contact.Data.AgeGroup,
				})
		}
	}

	if request.Roles != nil {
		for _, role := range request.Roles.Remove {
			contact.Data.Roles = slice.RemoveInPlace(role, contact.Data.Roles)
		}
		contact.Data.Roles = append(contact.Data.Roles, request.Roles.Add...)
		updatedContactFields = append(updatedContactFields, "roles")
		contactUpdates = append(contactUpdates, dal.Update{Field: "roles", Value: contact.Data.Roles})
		params.TeamModuleUpdates = append(params.TeamModuleUpdates,
			dal.Update{
				Field: fmt.Sprintf("contacts.%s.roles", request.ContactID),
				Value: contact.Data.Roles,
			})
	}

	if request.RelatedTo != nil {
		recordRef := models4linkage.TeamModuleDocRef{
			ModuleID:   const4contactus.ModuleID,
			Collection: const4contactus.ContactsCollection,
			TeamID:     request.TeamID,
			ItemID:     request.ContactID,
		}
		relatableAdapted := facade4linkage.NewRelatableAdapter[*models4contactus.ContactDto](func(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (err error) {
			// Verify contactID belongs to the same team
			teamContactBriefID := recordRef.ItemID
			if _, existingContact := params.TeamModuleEntry.Data.Contacts[teamContactBriefID]; !existingContact {
				if _, err = GetContactByID(ctx, tx, params.Team.ID, recordRef.ItemID); err != nil {
					return fmt.Errorf("failed to get related contact: %w", err)
				}
			}
			return nil
		})
		var relUpdate []dal.Update
		if relUpdate, err = facade4linkage.SetRelated(ctx, tx, params.UserID, params.Started, relatableAdapted, params.Contact, recordRef, *request.RelatedTo); err != nil {
			return err
		}
		contactUpdates = append(contactUpdates, relUpdate...)
	}

	if len(contactUpdates) > 0 {
		contact.Data.IncreaseVersion(params.Started, params.UserID)
		contactUpdates = append(contactUpdates, contact.Data.WithUpdatedAndVersion.GetUpdates()...)
		if err := contact.Data.Validate(); err != nil {
			return fmt.Errorf("contact DTO is not valid after updating %d fields (%+v) and before storing changes DB: %w",
				len(updatedContactFields), updatedContactFields, err)
		}
		if err := tx.Update(ctx, contact.Key, contactUpdates); err != nil {
			return err
		}
	}

	return nil
}
