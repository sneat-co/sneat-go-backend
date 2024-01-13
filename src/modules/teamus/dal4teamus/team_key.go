package dal4teamus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core"
)

// TeamsCollection table name
const TeamsCollection = "teams"

//const TeamBriefsCollection = "briefs"

// NewTeamKey create new doc ref
func NewTeamKey(id string) *dal.Key {
	const maxLen = 30
	if id == "" {
		panic("empty team ID")
	}
	if l := len(id); l > maxLen {
		panic(fmt.Sprintf("team ID is %v characters long exceded what is %v more then maxLen %v", l, maxLen-l, maxLen))
	}
	if !core.IsAlphanumericOrUnderscore(id) {
		panic(fmt.Sprintf("team ID has non alphanumeric characters or letters in upper case: [%v]", id))
	}
	return dal.NewKeyWithID(TeamsCollection, id)
}
