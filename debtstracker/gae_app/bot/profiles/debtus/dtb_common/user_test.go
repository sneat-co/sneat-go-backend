package dtb_common

import "testing"

func TestGetUserWithNilContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	if _, err := GetUser(nil); err != nil {
		t.Error("unexpected error", err)
	}
}
