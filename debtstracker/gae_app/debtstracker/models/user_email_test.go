package models

import (
	"testing"

	"github.com/strongo/strongoapp/appuser"
)

func TestUserEmail(t *testing.T) {
	var _ appuser.AccountData = (*UserEmailData)(nil)
}

func TestUserEmailEntity(t *testing.T) {
	var _ appuser.AccountData = (*UserEmailData)(nil)
}

func TestUserEmailEntity_AddProvider(t *testing.T) {
	entity := new(UserEmailData)

	if changed := entity.AddProvider("facebook"); !changed {
		t.Error("Should return changed=true")
	}
	if providerCount := len(entity.Providers); providerCount != 1 {
		t.Errorf("Expected to have 1 provider, got: %d", providerCount)
	}
	if changed := entity.AddProvider("facebook"); changed {
		t.Error("Should return changed=false")
	}
}
