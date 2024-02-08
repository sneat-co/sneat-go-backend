package gaedal

import (
	"github.com/dal-go/dalgo/dal"
	"testing"
)

func testStrKey(t *testing.T, expectedID string, key *dal.Key) {
	if key == nil {
		t.Error("key is nil")
		return
	}
	switch id := key.ID.(type) {
	case string:
		if id != expectedID {
			t.Error("IntegerID() != expectedID", expectedID)
		}
	case int, int64:
		t.Error("IntID() is not empty")
	}
	if key.Parent() != nil {
		t.Error("Parent() != nil")
	}
}

func testIntKey(t *testing.T, expectedID int64, key *dal.Key) {
	if key == nil {
		t.Error("key is nil")
		return
	}
	switch id := key.ID.(type) {
	case string:
		t.Error("StringID() is not empty")
	case int64:
		if id != expectedID {
			t.Error("IntegerID() != expectedID", expectedID)
		}
	case int:
		if id != int(expectedID) {
			t.Error("IntegerID() != expectedID", expectedID)
		}
	}
	if key.Parent() != nil {
		t.Error("Parent() != nil")
	}
}

func testDatastoreStringKey(t *testing.T, expectedID string, key *dal.Key) {
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

func testIncompleteKey(t *testing.T, key *dal.Key) {
	if key == nil {
		t.Error("key is nil")
		return
	}
	//if key.StringID() != "" {
	//	t.Error("StringID() is not empty")
	//}
	//if key.IntID() != 0 {
	//	t.Error("IntegerID() != 0")
	//}
	//if key.Parent() != nil {
	//	t.Error("Parent() != nil")
	//}
}
