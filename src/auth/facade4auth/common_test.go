package facade4auth

import (
	"github.com/dal-go/dalgo/dal"
	"testing"
)

func testStringKey(t *testing.T, expectedID string, key *dal.Key) {
	if key == nil {
		t.Error("key is nil")
		return
	}
	if id, ok := key.ID.(string); ok && id != expectedID {
		t.Error("StringID() != expectedID", id, expectedID)
	}
	if id, ok := key.ID.(int); ok && id != 0 {
		t.Error("IntegerID() != 0")
	}
	if key.Parent() != nil {
		t.Error("Parent() != nil")
	}
}
