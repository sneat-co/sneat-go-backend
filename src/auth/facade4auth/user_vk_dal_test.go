package facade4auth

import (
	"testing"
)

func TestNewVkUserKey(t *testing.T) {
	const vkUserID = 789
	key := NewUserVkKey(vkUserID)
	if key.ID.(string) != "789" {
		t.Error("key.ContactID != 789")
	}
}
