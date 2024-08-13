package models4auth

import (
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
)

const UserGooglePlusKind = "UserGooglePlus"

type UserGooglePlusEntity struct {
	appuser.OwnedByUserWithID
	Email          string `datastore:",noindex"`
	DisplayName    string `datastore:",noindex"`
	RefreshToken   string `datastore:",noindex"`
	ServerAuthCode string `datastore:",noindex"`
	AccessToken    string `datastore:",noindex"`
	ImageUrl       string `datastore:",noindex"`
	IdToken        string `datastore:",noindex"`

	Locale        string `datastore:",noindex"`
	NameFirst     string `datastore:",noindex"`
	NameLast      string `datastore:",noindex"`
	EmailVerified bool   `datastore:",noindex"`
}

type UserGooglePlus = record.DataWithID[string, *UserGooglePlusEntity]
