package facade4contactus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dto4contactus"
	"github.com/strongo/slice"
	"slices"
)

func updateContactRoles(params *dal4contactus.ContactWorkerParams, roles dto4contactus.SetRolesRequest) (updatedContactFields []string, err error) {
	var removedCount int
	var addedCount int
	params.Contact.Data.Roles, removedCount = slice.RemoveInPlace(params.Contact.Data.Roles, func(v string) bool {
		return slices.Contains(roles.Remove, v)
	})
	for _, role := range roles.Add {
		if !slices.Contains(params.Contact.Data.Roles, role) {
			addedCount++
			params.Contact.Data.Roles = append(params.Contact.Data.Roles, role)
		}
	}
	if removedCount > 0 || addedCount > 0 {
		updatedContactFields = append(updatedContactFields, "roles")
		params.ContactUpdates = append(params.ContactUpdates, dal.Update{Field: "roles", Value: params.Contact.Data.Roles})
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
			dal.Update{
				Field: fmt.Sprintf("contacts.%s.roles", params.Contact.ID),
				Value: params.Contact.Data.Roles,
			})
	}

	return updatedContactFields, err
}

func removeContactRoles(params *dal4contactus.ContactWorkerParams) {
	contact := params.Contact
	contactBrief := params.SpaceModuleEntry.Data.GetContactBriefByContactID(contact.ID)
	if contactBrief != nil {
		for _, update := range contactBrief.RemoveRole(const4contactus.SpaceMemberRoleMember) {
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, dal.Update{
				Field: fmt.Sprintf("contacts.%s.roles.%s", contact.ID, update.Field),
				Value: contact.Data.Roles,
			})
		}
	}
}
