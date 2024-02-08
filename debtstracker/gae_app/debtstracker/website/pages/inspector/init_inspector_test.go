package inspector

import "testing"

func TestInitInspector(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	InitInspector(nil)
}
