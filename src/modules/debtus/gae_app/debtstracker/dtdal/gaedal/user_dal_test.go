package gaedal

import (
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"testing"
)

func TestNewAppUserKey(t *testing.T) {
	const appUserID = "1234"
	testStrKey(t, appUserID, dbo4userus.NewUserKey(appUserID))
}
