package invites

import "testing"

func TestSendInviteByEmail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	_, _ = SendInviteByEmail(nil, nil, "", "", "", "", "", "")
}
