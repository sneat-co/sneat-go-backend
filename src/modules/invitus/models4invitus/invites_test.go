package models4invitus

import "testing"

func TestInvite_Validate(t *testing.T) {
	invite := InviteDto{}
	if err := invite.Validate(); err == nil {
		t.Fatal("error expected for empty value")
	}
}
