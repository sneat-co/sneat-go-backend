package api4userus

import "testing"

func TestSetUserTitle(t *testing.T) {
	defer func() {
		if p := recover(); p == nil {
			t.Fatal("panic expected")
		}
	}()
	httpInitUserRecord(nil, nil)
	// TODO: implement positive
}
