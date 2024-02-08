package admin

import "testing"

func TestInitAdmin(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	InitAdmin(nil)
}
