package dal4contactus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core"
)

// NewContactKey creates a new contact's key in format "spaceID:memberID"
func NewContactKey(teamID, contactID string) *dal.Key {
	if !core.IsAlphanumericOrUnderscore(contactID) {
		panic(fmt.Errorf("contactID should be alphanumeric, got: [%s]", contactID))
	}
	teamModuleKey := dbo4spaceus.NewSpaceModuleKey(teamID, const4contactus.ModuleID)
	return dal.NewKeyWithParentAndID(teamModuleKey, dbo4contactus.SpaceContactsCollection, contactID)
}
