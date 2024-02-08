package gaedal

import (
	"testing"
)

func TestNewGoogleUserKey(t *testing.T) {
	const googleUserID = "246"
	testDatastoreStringKey(t, googleUserID, NewUserGoogleKey(googleUserID))
}
