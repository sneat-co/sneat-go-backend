package common4all

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeID(t *testing.T) {
	_, err := DecodeID("")
	assert.NotNil(t, err) // Should return error if empty string
}
