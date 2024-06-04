package models4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
)

func NewContactFullRef(teamID, contactID string) dbo4linkage.TeamModuleItemRef {
	return dbo4linkage.NewTeamModuleItemRef(teamID, const4contactus.ModuleID, const4contactus.ContactsCollection, contactID)
}
