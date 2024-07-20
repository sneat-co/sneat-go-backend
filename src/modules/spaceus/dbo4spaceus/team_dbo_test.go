package dbo4spaceus

import (
	"testing"
)

func TestSpace_Validate(t *testing.T) {
	t.Run("must_fail", func(t *testing.T) {
		v := SpaceDbo{}
		if err := v.Validate(); err == nil {
			t.Fatal("Expected to get validation error for empty team record")
		}
	})
}
