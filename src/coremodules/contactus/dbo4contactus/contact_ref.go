package dbo4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/linkage/dbo4linkage"
)

func NewContactFullRef(teamID, contactID string) dbo4linkage.SpaceModuleItemRef {
	return dbo4linkage.NewSpaceModuleItemRef(teamID, const4contactus.ModuleID, const4contactus.ContactsCollection, contactID)
}
