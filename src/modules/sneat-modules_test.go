package modules

import (
	"testing"
)

func TestModules(t *testing.T) {
	modules := Modules()
	if len(modules) == 0 {
		t.Error("len(modules) == 0")
	}
}
