package models4teamus

import (
	"testing"
)

func TestTeam_Validate(t *testing.T) {
	t.Run("must_fail", func(t *testing.T) {
		v := TeamDto{}
		if err := v.Validate(); err == nil {
			t.Fatal("Expected to get validation error for empty team record")
		}
	})
}
