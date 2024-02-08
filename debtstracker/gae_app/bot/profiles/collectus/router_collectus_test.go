package collectus

import "testing"

func TestRouter(t *testing.T) {
	if Router.CommandsCount() == 0 {
		t.Fatal("Router.CommandsCount() == 0")
	}
}
