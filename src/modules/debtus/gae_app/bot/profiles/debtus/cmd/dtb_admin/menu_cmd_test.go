package dtb_admin

import "testing"

func TestAdmin(t *testing.T) {
	if AdminCommand.Code == "" {
		t.Fatal("AdminCommand.Code is not set")
	}
}
