package support

import "testing"

func TestInitSupportHandlers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	InitSupportHandlers(nil)
}
