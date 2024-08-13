package models4auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserGoogleEntity_GetEmail(t *testing.T) {

	userAccount := NewUserAccountEntry("1")
	userAccount.Data.EmailLowerCase = "test@example.com"
	assert.Equal(t, "test@example.com", userAccount.Data.GetEmailLowerCase())
}
