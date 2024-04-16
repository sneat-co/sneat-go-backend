package facade4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
)

func removeContactRoles(
	params *dal4contactus.ContactusTeamWorkerParams,
	contact dal4contactus.ContactEntry,
) (contactUpdates []dal.Update) {
	contactBrief := params.TeamModuleEntry.Data.GetContactBriefByContactID(contact.ID)
	if contactBrief != nil && contactBrief.RemoveRole(const4contactus.TeamMemberRoleMember) {
		params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{Field: "contacts." + contact.ID + ".roles", Value: contact.Data.Roles})
	}

	if contact.Data.RolesField.RemoveRole(const4contactus.TeamMemberRoleMember) {
		contactUpdates = append(contactUpdates, dal.Update{Field: "roles", Value: contact.Data.Roles})
	}
	return contactUpdates
}
