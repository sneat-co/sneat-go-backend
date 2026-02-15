package dal4listus

import (
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"testing"
)

func TestNewListKey(t *testing.T) {
	spaceID := coretypes.SpaceID("s1")
	listKey := dbo4listus.ListKey("do!123")
	key := NewListKey(spaceID, listKey)
	if key.ID != "do!123" {
		t.Errorf("key.ID = %v, want do!123", key.ID)
	}
	if key.Collection() != dbo4listus.ListsCollection {
		t.Errorf("key.Collection() = %v, want %v", key.Collection(), dbo4listus.ListsCollection)
	}
}

func TestNewListEntry(t *testing.T) {
	spaceID := coretypes.SpaceID("s1")
	listKey := dbo4listus.ListKey("do!123")
	entry := NewListEntry(spaceID, listKey)
	if entry.ID != "do!123" {
		t.Errorf("entry.ID = %v, want do!123", entry.ID)
	}
	if entry.Data == nil {
		t.Error("entry.Data is nil")
	}
}
