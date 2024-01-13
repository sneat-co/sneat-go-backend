package models4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/models4invitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/strongo/strongoapp/with"
)

// TeamContactsCollection defines  collection name for team contacts.
// We have `Team` prefix as it can belong only to single team
// and TeamID is also in record key as prefix.
const TeamContactsCollection = "contacts"

// ContactDto belongs only to single team
type ContactDto struct {
	//dbmodels.WithTeamID -- not needed as it's in record key
	//dbmodels.WithUserIDs

	models4linkage.WithRelatedAndIDs
	briefs4contactus.ContactBase
	with.CreatedFields
	with.TagsField
	briefs4contactus.WithMultiTeamContacts[*briefs4contactus.ContactBrief]
	models4invitus.WithInvites // Invites to become a team member
}

// Validate returns error if not valid
func (v ContactDto) Validate() error {
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
