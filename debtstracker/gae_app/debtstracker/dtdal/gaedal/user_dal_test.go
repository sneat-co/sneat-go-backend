package gaedal

import (
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"testing"
)

func TestNewAppUserKey(t *testing.T) {
	const appUserID = "1234"
	testStrKey(t, appUserID, models.NewAppUserKey(appUserID))
}
