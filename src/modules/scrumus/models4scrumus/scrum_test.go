package models4scrumus

import (
	"testing"
)

func TestScrum_Validate(t *testing.T) {
	t.Run("empty_record", func(t *testing.T) {
		record := Scrum{}
		if err := record.Validate(); err != nil {
			t.Fatalf("no error expected for empty value, got: %v", err)
		}
	})
}
