package splitus

import "testing"

func TestGetBillID(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	_, _ = GetBillID(nil)
}
