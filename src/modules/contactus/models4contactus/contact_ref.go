package models4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
)

func NewContactFullRef(teamID, contactID string) models4linkage.TeamModuleItemRef {
	return models4linkage.NewTeamModuleDocRef(teamID, const4contactus.ModuleID, const4contactus.ContactsCollection, contactID)
}
