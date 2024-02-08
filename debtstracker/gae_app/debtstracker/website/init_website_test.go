package website

import "testing"

func TestInitWebsite(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	InitWebsite(nil)
}
