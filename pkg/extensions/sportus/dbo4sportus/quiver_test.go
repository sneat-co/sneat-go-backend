package dbo4sportus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel_Validate(t *testing.T) {
	// Model.Validate currently has no constraints and always succeeds.
	assert.NoError(t, Model{}.Validate())
	assert.NoError(t, Model{Brand: "Duotone", Title: "Neo SLS", Kinds: []string{"kiting"}}.Validate())
}
