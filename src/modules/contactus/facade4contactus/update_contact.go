package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-core/facade"
)

// UpdateContact sets contact fields
func UpdateContact(
	ctx context.Context,
	user facade.User,
	request dto4contactus.UpdateContactRequest,
) (err error) {
	return dal4contactus.RunContactWorker(ctx, user, request.ContactRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
			return UpdateContactTx(ctx, tx, request, params)
		})
}

// UpdateContactTx sets contact fields
func UpdateContactTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4contactus.UpdateContactRequest,
	params *dal4contactus.ContactWorkerParams,
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
	params *dal4contactus.ContactWorkerParams,
) (err error) {

	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}

	contact := params.Contact

	if err := contact.Data.Validate(); err != nil {
		return fmt.Errorf("contact DTO is not valid after loading from DB: %w", err)
	}

	contactBrief := params.TeamModuleEntry.Data.Contacts[request.ContactID]

	var updatedContactFields []string

	if request.Address != nil {
		if *request.Address != *contact.Data.Address {
			updatedContactFields = append(updatedContactFields, "address")
			contact.Data.Address = request.Address
			params.ContactUpdates = append(params.ContactUpdates, dal.Update{Field: "address", Value: request.Address})
		}
	}

	if request.VatNumber != nil {
		if vat := *request.VatNumber; vat != contact.Data.VATNumber {
			updatedContactFields = append(updatedContactFields, "vatNumber")
			contact.Data.VATNumber = vat
			params.ContactUpdates = append(params.ContactUpdates, dal.Update{Field: "vatNumber", Value: vat})
		}
	}

	if request.AgeGroup != "" {
		if request.AgeGroup != contact.Data.AgeGroup {
			updatedContactFields = append(updatedContactFields, "ageGroup")
			contact.Data.AgeGroup = request.AgeGroup
			params.ContactUpdates = append(params.ContactUpdates, dal.Update{Field: "ageGroup", Value: contact.Data.AgeGroup})
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
		var contactFieldsUpdated []string
		if contactFieldsUpdated, err = updateContactRoles(params, *request.Roles); err != nil {
			return err
		}
		updatedContactFields = append(updatedContactFields, contactFieldsUpdated...)
	}

	if request.Related != nil {
		itemRef := models4linkage.TeamModuleItemRef{
			ModuleID:   const4contactus.ModuleID,
			Collection: const4contactus.ContactsCollection,
			TeamID:     request.TeamID,
			ItemID:     request.ContactID,
		}
		if err = updateRelatedField(ctx, tx, itemRef, request.UpdateRelatedRequest, params); err != nil {
			return err
		}
	}

	if len(params.ContactUpdates) > 0 {
		contact.Data.IncreaseVersion(params.Started, params.UserID)
		params.ContactUpdates = append(params.ContactUpdates, contact.Data.WithUpdatedAndVersion.GetUpdates()...)
		if err := contact.Data.Validate(); err != nil {
			return fmt.Errorf("contact DTO is not valid after updating %d fields (%+v) and before storing changes DB: %w",
				len(updatedContactFields), updatedContactFields, err)
		}
		if err := tx.Update(ctx, contact.Key, params.ContactUpdates); err != nil {
			return err
		}
	}

	return nil
}
