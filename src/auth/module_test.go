package auth

import (
	"testing"
)

func TestModule(t *testing.T) {
	m := Module()
	if m == nil {
		t.Fatal("Module() returned nil")
	}
}
