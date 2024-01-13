package sneatbot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	assert.NotNil(t, Profile)
}
