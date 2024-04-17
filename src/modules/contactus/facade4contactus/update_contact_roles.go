package facade4contactus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/strongo/slice"
)

func updateContactRoles(params *dal4contactus.ContactWorkerParams, roles dto4contactus.SetRolesRequest) (updatedContactFields []string, err error) {
	var rolesUpdated bool
	for _, role := range roles.Remove {
		roles := slice.RemoveInPlace(role, params.Contact.Data.Roles)
		rolesUpdated = rolesUpdated || len(roles) != len(params.Contact.Data.Roles)
		params.Contact.Data.Roles = roles
	}

	for _, role := range roles.Add {
		if !slice.Contains(params.Contact.Data.Roles, role) {
			rolesUpdated = true
			params.Contact.Data.Roles = append(params.Contact.Data.Roles, role)
		}
	}
	if rolesUpdated {
		updatedContactFields = append(updatedContactFields, "roles")
		params.ContactUpdates = append(params.ContactUpdates, dal.Update{Field: "roles", Value: params.Contact.Data.Roles})
		params.TeamModuleUpdates = append(params.TeamModuleUpdates,
			dal.Update{
				Field: fmt.Sprintf("contacts.%s.roles", params.Contact.ID),
				Value: params.Contact.Data.Roles,
			})
	}

	return updatedContactFields, err
}

func removeContactRoles(
	params *dal4contactus.ContactWorkerParams,
) {
	contact := params.Contact
	contactBrief := params.TeamModuleEntry.Data.GetContactBriefByContactID(contact.ID)
	if contactBrief != nil && contactBrief.RemoveRole(const4contactus.TeamMemberRoleMember) {
		params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{Field: "contacts." + contact.ID + ".roles", Value: contact.Data.Roles})
	}

	if contact.Data.RolesField.RemoveRole(const4contactus.TeamMemberRoleMember) {
		params.ContactUpdates = append(params.ContactUpdates, dal.Update{Field: "roles", Value: contact.Data.Roles})
	}
}
