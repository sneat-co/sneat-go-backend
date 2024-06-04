package dbo4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
)

// UsersCollection a name of the user's db table
const UsersCollection = "users"

// UserDefaults keeps user's defaults
type UserDefaults struct {
	ShortNames []briefs4contactus.ShortName `json:"shortNames,omitempty" firestore:"shortNames,omitempty"`
}

type User = record.DataWithID[string, *UserDbo]

func NewUser(id string) User {
	key := dal.NewKeyWithID(UsersCollection, id)
	return record.NewDataWithID(id, key, new(UserDbo))
}
