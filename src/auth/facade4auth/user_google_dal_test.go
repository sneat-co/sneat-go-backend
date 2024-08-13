package facade4auth

import (
	"testing"
)

func TestNewGoogleUserKey(t *testing.T) {
	const googleUserID = "246"
	testStringKey(t, googleUserID, NewUserGoogleKey(googleUserID))
}
