package models4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/strongo/strongoapp/with"
)

// SpaceContactsCollection defines collection name for team contacts.
// We have `Space` prefix as it can belong only to a single team
// and SpaceID is also in record key as prefix.
const SpaceContactsCollection = "contacts"

// ContactDbo belongs only to a single team
type ContactDbo struct {
	//dbmodels.WithSpaceID -- not needed as it's in record key
	//dbmodels.WithUserIDs

	briefs4contactus.ContactBase

	dbo4linkage.WithRelatedAndIDs
	with.CreatedFields
	with.TagsField
	briefs4contactus.WithMultiSpaceContacts[*briefs4contactus.ContactBrief]
	dbo4invitus.WithInvites // Invites to become a team member
}

// Validate returns error if not valid
func (v ContactDbo) Validate() error {
	if err := v.ContactBase.Validate(); err != nil {
		return fmt.Errorf("ContactRecordBase is not valid: %w", err)
	}
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	if err := v.RolesField.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	if err := v.WithInvites.Validate(); err != nil {
		return err
	}
	if err := v.WithRelatedAndIDs.Validate(); err != nil {
		return err
	}
	return nil
}
