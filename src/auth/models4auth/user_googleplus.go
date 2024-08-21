package models4auth

import (
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
)

const UserGooglePlusKind = "UserGooglePlus"

type UserGooglePlusEntity struct {
	appuser.OwnedByUserWithID
	Email          string `firestore:",omitempty"`
	DisplayName    string `firestore:",omitempty"`
	RefreshToken   string `firestore:",omitempty"`
	ServerAuthCode string `firestore:",omitempty"`
	AccessToken    string `firestore:",omitempty"`
	ImageUrl       string `firestore:",omitempty"`
	IdToken        string `firestore:",omitempty"`

	Locale        string `firestore:",omitempty"`
	NameFirst     string `firestore:",omitempty"`
	NameLast      string `firestore:",omitempty"`
	EmailVerified bool   `firestore:",omitempty"`
}

type UserGooglePlus = record.DataWithID[string, *UserGooglePlusEntity]
