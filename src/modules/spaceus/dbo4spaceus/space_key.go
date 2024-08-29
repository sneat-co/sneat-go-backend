package dbo4spaceus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core"
)

// SpacesCollection table name
const SpacesCollection = "spaces"
const SpacesFiled = "spaces"

//const SpaceBriefsCollection = "briefs"

// NewSpaceKey create new doc ref
func NewSpaceKey(spaceID string) *dal.Key {
	const maxLen = 30
	if spaceID == "" {
		panic("spaceID is empty string")
	}
	if l := len(spaceID); l > maxLen {
		panic(fmt.Sprintf("space ContactID is %v characters long exceded what is %d more then maxLen %d", l, maxLen-l, maxLen))
	}
	if !core.IsAlphanumericOrUnderscore(spaceID) {
		panic(fmt.Sprintf("space ContactID has non alphanumeric characters or letters in upper case: [%s]", spaceID))
	}
	return dal.NewKeyWithID(SpacesCollection, spaceID)
}
