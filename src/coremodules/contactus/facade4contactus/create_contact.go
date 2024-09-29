package facade4contactus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/person"
	"slices"
)

// CreateContact creates team contact
func CreateContact(
	ctx context.Context,
	userCtx facade.UserContext,
	userCanBeNonSpaceMember bool,
	request dto4contactus.CreateContactRequest,
) (
	response dto4contactus.CreateContactResponse,
	err error,
) {
	if err = request.Validate(); err != nil {
		return response, fmt.Errorf("invalid CreateContactRequest: %w", err)
	}

	err = dal4spaceus.CreateSpaceItem(ctx, userCtx, request.SpaceRequest, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4contactus.ContactusSpaceDbo]) (err error) {
			var contact dal4contactus.ContactEntry
			if contact, err = CreateContactTx(ctx, tx, userCanBeNonSpaceMember, request, params); err != nil {
				return err
			}
			response = dto4contactus.CreateContactResponse{
				ID:   contact.ID,
				Data: contact.Data,
			}
			if response.Data == nil {
				return errors.New("CreateContactTx returned nil contact data")
			}
			return err
		},
	)
	if err != nil {
		err = fmt.Errorf("failed to create a new contact: %w", err)
		return
	}
	return
}

func CreateContactTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCanBeNonSpaceMember bool,
	request dto4contactus.CreateContactRequest,
	params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4contactus.ContactusSpaceDbo],
) (
	contact dal4contactus.ContactEntry,
	err error,
) {
	if err = request.Validate(); err != nil {
		return
	}
	if err = params.GetRecords(ctx, tx); err != nil {
		return
	}
	userContactID, userContactBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(params.UserID())
	if !userCanBeNonSpaceMember && (userContactBrief == nil || !userContactBrief.IsSpaceMember()) {
		err = errors.New("user is not a member of the team")
		return
	}
	if request.Related != nil {
		relatedByCollection := request.Related[const4contactus.ModuleID]
		if relatedByCollection != nil {
			relatedItems := relatedByCollection[const4contactus.ContactsCollection]
			if len(relatedItems) > 0 {
				var isRelatedByUserID bool
				for _, relatedItem := range relatedItems {
					isRelatedByUserID = dbo4linkage.HasRelatedItem(relatedItems, dbo4linkage.RelatedItemKey{SpaceID: params.Space.ID, ItemID: params.UserID()})
					if !isRelatedByUserID {
						contactID := relatedItem.Keys[0].ItemID
						if contactBrief := params.SpaceModuleEntry.Data.GetContactBriefByContactID(contactID); contactBrief == nil {
							return contact, fmt.Errorf("contact with ContactID=[%s] is not found", contactID)
						}
					}
					switch userContactBrief.AgeGroup {
					case "", dbmodels.AgeGroupUnknown:
						for relatedAs := range relatedItem.RolesOfItem {
							switch relatedAs {
							case dbmodels.RelationshipSpouse, dbmodels.RelationshipChild:
								userContactBrief.AgeGroup = dbmodels.AgeGroupAdult
								userContactKey := dal4contactus.NewContactKey(request.SpaceID, userContactID)
								if err = tx.Update(ctx, userContactKey, []dal.Update{
									{
										Field: "ageGroup",
										Value: userContactBrief.AgeGroup,
									},
								}); err != nil {
									err = fmt.Errorf("failed to update member record: %w", err)
									return
								}
							}
						}
					}
				}
				if isRelatedByUserID {
					userRelatedItem := dbo4linkage.GetRelatedItemByKey(relatedItems, dbo4linkage.RelatedItemKey{SpaceID: params.Space.ID, ItemID: params.UserID()})
					userRelatedItem.Keys[0].ItemID = userContactID
				}
			}
		}
	}

	parentContactID := request.ParentContactID

	var parent dal4contactus.ContactEntry
	if parentContactID != "" {
		parent = dal4contactus.NewContactEntry(request.SpaceID, parentContactID)
		if err = tx.Get(ctx, parent.Record); err != nil {
			return contact, fmt.Errorf("failed to get parent contact with ContactID=[%s]: %w", parentContactID, err)
		}
	}

	contactDbo := new(dbo4contactus.ContactDbo)
	contactDbo.CreatedAt = params.Started
	contactDbo.CreatedBy = params.UserID()
	contactDbo.Status = "active"
	contactDbo.ParentID = parentContactID
	contactDbo.RolesField = request.RolesField
	if request.Person != nil {
		contactDbo.ContactBase = request.Person.ContactBase
		contactDbo.Type = briefs4contactus.ContactTypePerson
		if contactDbo.AgeGroup == "" {
			contactDbo.AgeGroup = "unknown"
		}
		if contactDbo.Gender == "" {
			contactDbo.Gender = "unknown"
		}
		contactDbo.ContactBase = request.Person.ContactBase
		for _, role := range request.Roles {
			if !slices.Contains(contactDbo.Roles, role) {
				contactDbo.Roles = append(contactDbo.Roles, role)
			}
		}
	} else if request.Company != nil {
		contactDbo.Type = briefs4contactus.ContactTypeCompany
		contactDbo.Title = request.Company.Title
		contactDbo.VATNumber = request.Company.VATNumber
		contactDbo.Address = request.Company.Address
	} else if request.Location != nil {
		contactDbo.Type = briefs4contactus.ContactTypeLocation
		contactDbo.Title = request.Location.Title
		contactDbo.Address = &request.Location.Address
	} else if request.Basic != nil {
		contactDbo.Type = request.Type
		contactDbo.Title = request.Basic.Title
	} else {
		return contact, errors.New("contact type is not specified")
	}
	if contactDbo.Address != nil {
		contactDbo.CountryID = contactDbo.Address.CountryID
	}
	contactDbo.ShortTitle = contactDbo.DetermineShortTitle(request.Person.Title, params.SpaceModuleEntry.Data.Contacts)
	var contactID string
	if request.ContactID == "" {
		contactIDs := params.SpaceModuleEntry.Data.ContactIDs()
		if contactID, err = person.GenerateIDFromNameOrRandom(request.Person.Names, contactIDs); err != nil {
			return contact, fmt.Errorf("failed to generate contact ContactID: %w", err)
		}
	} else {
		contactID = request.ContactID
	}
	if contactDbo.CountryID == "" && params.Space.Data.CountryID != "" && params.Space.Data.Type == core4spaceus.SpaceTypeFamily {
		contactDbo.CountryID = params.Space.Data.CountryID
	}
	params.SpaceModuleEntry.Data.AddContact(contactID, &contactDbo.ContactBrief)
	if params.SpaceModuleEntry.Record.Exists() {
		if err = tx.Update(ctx, params.SpaceModuleEntry.Key, []dal.Update{
			{
				Field: const4contactus.ContactsField,
				Value: params.SpaceModuleEntry.Data.Contacts,
			},
		}); err != nil {
			return contact, fmt.Errorf("failed to update team contact briefs: %w", err)
		}
	} else {
		if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
			return contact, fmt.Errorf("faield to insert team contacts brief record: %w", err)
		}
	}

	//params.SpaceUpdates = append(params.SpaceUpdates, params.Space.Data.UpdateNumberOf(const4contactus.ContactsField, len(params.SpaceModuleEntry.Data.Contacts)))

	if request.Related != nil {
		if err = updateRelationshipsInRelatedItems(ctx, tx, params.UserID(), userContactID, params.Space.ID, contactID, params.SpaceModuleEntry, contactDbo, request.Related); err != nil {
			err = fmt.Errorf("failed to update relationships in related items: %w", err)
			return
		}
	}

	contact = dal4contactus.NewContactEntryWithData(request.SpaceID, contactID, contactDbo)

	_ = dbo4linkage.UpdateRelatedIDs(&contact.Data.WithRelated, &contact.Data.WithRelatedIDs)
	if err = contact.Data.Validate(); err != nil {
		return contact, fmt.Errorf("contact record is not valid: %w", err)
	}
	if err = tx.Insert(ctx, contact.Record); err != nil {
		return contact, fmt.Errorf("failed to insert contact record: %w", err)
	}
	if parent.ID != "" {
		if err = updateParentContact(ctx, tx, contact, parent); err != nil {
			return contact, fmt.Errorf("failed to update parent contact: %w", err)
		}
	}
	return contact, err
}

func updateRelationshipsInRelatedItems(ctx context.Context, tx dal.ReadTransaction,
	userID, userContactID, spaceID, contactID string,
	contactusSpaceEntry dal4contactus.ContactusSpaceEntry,
	contactDbo *dbo4contactus.ContactDbo,
	related dbo4linkage.RelatedByModuleID,
) (err error) {
	if userContactID == "" { // Why we get it 2nd time? Previous is up in stack in CreateContactTx()
		if userContactID, err = dal4userus.GetUserSpaceContactID(ctx, tx, userID, contactusSpaceEntry); err != nil {
			return
		}
		if userContactID == "" {
			err = errors.New("user is not associated with the spaceID=" + spaceID)
			return
		}
	}

	for moduleID, relatedByCollection := range related {
		for collection, relatedByItemID := range relatedByCollection {
			for _, relatedItem := range relatedByItemID {
				for _, key := range relatedItem.Keys {
					itemRef := dbo4linkage.SpaceModuleItemRef{
						Space:      spaceID,
						Module:     moduleID,
						Collection: collection,
						ItemID:     key.ItemID,
					}

					if _, err = contactDbo.AddRelationshipsAndIDs(
						itemRef,
						relatedItem.RolesOfItem,
						relatedItem.RolesToItem,
					); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
