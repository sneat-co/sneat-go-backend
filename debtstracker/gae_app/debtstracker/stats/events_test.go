package stats

import "testing"

func TestConstants(t *testing.T) {
	if USER_CREATED != "USER_CREATED" {
		t.Fatal("Missing event USER_CREATED")
	}
}
