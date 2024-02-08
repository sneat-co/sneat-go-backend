package vkapp

import "testing"

func TestInitVkIFrameAppWithoutRouter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	InitVkIFrameApp(nil)
}
