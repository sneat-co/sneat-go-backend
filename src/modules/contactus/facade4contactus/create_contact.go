package facade4contactus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/person"
)

// CreateContact creates team contact
func CreateContact(
	ctx context.Context,
	userContext facade.User,
	userCanBeNonTeamMember bool,
	request dto4contactus.CreateContactRequest,
) (
	response dto4contactus.CreateContactResponse,
	err error,
) {
	if err = request.Validate(); err != nil {
		return response, fmt.Errorf("invalid CreateContactRequest: %w", err)
	}

	err = dal4teamus.CreateTeamItem(ctx, userContext, request.TeamRequest, const4contactus.ModuleID, new(models4contactus.ContactusTeamDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4contactus.ContactusTeamDbo]) (err error) {
			var contact dal4contactus.ContactEntry
			if contact, err = CreateContactTx(ctx, tx, userCanBeNonTeamMember, request, params); err != nil {
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
	userCanBeNonTeamMember bool,
	request dto4contactus.CreateContactRequest,
	params *dal4teamus.ModuleTeamWorkerParams[*models4contactus.ContactusTeamDbo],
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
	userContactID, userContactBrief := params.TeamModuleEntry.Data.GetContactBriefByUserID(params.UserID)
	if !userCanBeNonTeamMember && (userContactBrief == nil || !userContactBrief.IsTeamMember()) {
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
					isRelatedByUserID = dbo4linkage.HasRelatedItem(relatedItems, dbo4linkage.RelatedItemKey{TeamID: params.Team.ID, ItemID: params.UserID})
					if !isRelatedByUserID {
						contactID := relatedItem.Keys[0].ItemID
						if contactBrief := params.TeamModuleEntry.Data.GetContactBriefByContactID(contactID); contactBrief == nil {
							return contact, fmt.Errorf("contact with ID=[%s] is not found", contactID)
						}
					}
					switch userContactBrief.AgeGroup {
					case "", dbmodels.AgeGroupUnknown:
						for relatedAs := range relatedItem.RolesOfItem {
							switch relatedAs {
							case dbmodels.RelationshipSpouse, dbmodels.RelationshipChild:
								userContactBrief.AgeGroup = dbmodels.AgeGroupAdult
								userContactKey := dal4contactus.NewContactKey(request.TeamID, userContactID)
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
					userRelatedItem := dbo4linkage.GetRelatedItemByKey(relatedItems, dbo4linkage.RelatedItemKey{TeamID: params.Team.ID, ItemID: params.UserID})
					userRelatedItem.Keys[0].ItemID = userContactID
				}
			}
		}
	}

	parentContactID := request.ParentContactID

	var parent dal4contactus.ContactEntry
	if parentContactID != "" {
		parent = dal4contactus.NewContactEntry(request.TeamID, parentContactID)
		if err = tx.Get(ctx, parent.Record); err != nil {
			return contact, fmt.Errorf("failed to get parent contact with ID=[%s]: %w", parentContactID, err)
		}
	}

	contactDbo := new(models4contactus.ContactDbo)
	contactDbo.CreatedAt = params.Started
	contactDbo.CreatedBy = params.UserID
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
			if !slice.Contains(contactDbo.Roles, role) {
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
	contactDbo.ShortTitle = contactDbo.DetermineShortTitle(request.Person.Title, params.TeamModuleEntry.Data.Contacts)
	var contactID string
	if request.ContactID == "" {
		contactIDs := params.TeamModuleEntry.Data.ContactIDs()
		if contactID, err = person.GenerateIDFromNameOrRandom(request.Person.Names, contactIDs); err != nil {
			return contact, fmt.Errorf("failed to generate contact ID: %w", err)
		}
	} else {
		contactID = request.ContactID
	}
	if contactDbo.CountryID == "" && params.Team.Data.CountryID != "" && params.Team.Data.Type == core4teamus.TeamTypeFamily {
		contactDbo.CountryID = params.Team.Data.CountryID
	}
	params.TeamModuleEntry.Data.AddContact(contactID, &contactDbo.ContactBrief)
	if params.TeamModuleEntry.Record.Exists() {
		if err = tx.Update(ctx, params.TeamModuleEntry.Key, []dal.Update{
			{
				Field: const4contactus.ContactsField,
				Value: params.TeamModuleEntry.Data.Contacts,
			},
		}); err != nil {
			return contact, fmt.Errorf("failed to update team contact briefs: %w", err)
		}
	} else {
		if err = tx.Insert(ctx, params.TeamModuleEntry.Record); err != nil {
			return contact, fmt.Errorf("faield to insert team contacts brief record: %w", err)
		}
	}

	//params.TeamUpdates = append(params.TeamUpdates, params.Team.Data.UpdateNumberOf(const4contactus.ContactsField, len(params.TeamModuleEntry.Data.Contacts)))

	if request.Related != nil {
		if err = updateRelationshipsInRelatedItems(ctx, tx, params.UserID, userContactID, params.Team.ID, contactID, params.TeamModuleEntry, contactDbo, request.Related); err != nil {
			err = fmt.Errorf("failed to update relationships in related items: %w", err)
			return
		}
	}

	contact = dal4contactus.NewContactEntryWithData(request.TeamID, contactID, contactDbo)

	_ = contact.Data.UpdateRelatedIDs()
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
	userID, userContactID, teamID, contactID string,
	contactusTeamEntry dal4contactus.ContactusTeamModuleEntry,
	contactDbo *models4contactus.ContactDbo,
	related dbo4linkage.RelatedByModuleID,
) (err error) {
	if userContactID == "" { // Why we get it 2nd time? Previous is up in stack in CreateContactTx()
		if userContactID, err = facade4userus.GetUserTeamContactID(ctx, tx, userID, contactusTeamEntry); err != nil {
			return
		}
		if userContactID == "" {
			err = errors.New("user is not associated with the teamID=" + teamID)
			return
		}
	}

	for moduleID, relatedByCollection := range related {
		for collection, relatedByItemID := range relatedByCollection {
			for _, relatedItem := range relatedByItemID {
				for _, key := range relatedItem.Keys {
					itemRef := dbo4linkage.TeamModuleItemRef{
						TeamID:     teamID,
						ModuleID:   moduleID,
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
