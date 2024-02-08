package models

import (
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/appuser"
	"testing"
)

func TestUserGoogleEntity_GetEmail(t *testing.T) {

	entity := UserAccount{
		data: &appuser.AccountDataBase{},
	}
	entity.data.EmailLowerCase = "test@example.com"
	assert.Equal(t, "test@example.com", entity.Data().GetEmailLowerCase())
}
