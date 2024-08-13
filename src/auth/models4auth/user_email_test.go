package models4auth

import (
	"testing"

	"github.com/strongo/strongoapp/appuser"
)

func TestUserEmail(t *testing.T) {
	var _ appuser.AccountData = (*UserEmailDbo)(nil)
}

func TestUserEmailEntity(t *testing.T) {
	var _ appuser.AccountData = (*UserEmailDbo)(nil)
}

func TestUserEmailEntity_AddProvider(t *testing.T) {
	entity := new(UserEmailDbo)

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
