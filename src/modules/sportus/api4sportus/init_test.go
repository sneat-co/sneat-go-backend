package api4sportus

import (
	"testing"
)

func TestRegisterRoutes(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("RegisterRoutes(nil) should panic")
		}
	}()
	RegisterRoutes(nil)
}
