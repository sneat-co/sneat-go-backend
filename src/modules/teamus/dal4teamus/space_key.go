package dal4teamus

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
func NewSpaceKey(id string) *dal.Key {
	const maxLen = 30
	if id == "" {
		panic("empty space ID")
	}
	if l := len(id); l > maxLen {
		panic(fmt.Sprintf("space ID is %v characters long exceded what is %d more then maxLen %d", l, maxLen-l, maxLen))
	}
	if !core.IsAlphanumericOrUnderscore(id) {
		panic(fmt.Sprintf("space ID has non alphanumeric characters or letters in upper case: [%s]", id))
	}
	return dal.NewKeyWithID(SpacesCollection, id)
}
