package models4logist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestContactKey tests getContactKey
func TestContactKey(t *testing.T) {
	assert.Equal(t, "contact=abc", getContactKey("abc"))
}

// TestIsContactKey tests isContactKey
func TestIsContactKey(t *testing.T) {
	assert.True(t, isContactKey(getContactKey("abc")))
	assert.False(t, isContactKey("contact-abc"))
}

// TestGetContactIdFromOrderKey tests getContactIdFromOrderKey
func TestGetContactIdFromOrderKey(t *testing.T) {
	const contactID = "abc"
	assert.Equal(t, contactID, getContactIdFromOrderKey(getContactKey(contactID)))
}
