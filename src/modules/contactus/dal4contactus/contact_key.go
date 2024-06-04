package dal4contactus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core"
)

// NewContactKey creates a new contact's key in format "teamID:memberID"
func NewContactKey(teamID, contactID string) *dal.Key {
	if !core.IsAlphanumericOrUnderscore(contactID) {
		panic(fmt.Errorf("contactID should be alphanumeric, got: [%s]", contactID))
	}
	teamModuleKey := dal4teamus.NewTeamModuleKey(teamID, const4contactus.ModuleID)
	return dal.NewKeyWithParentAndID(teamModuleKey, models4contactus.TeamContactsCollection, contactID)
}
