package dal4retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/const4retrospectus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"testing"
)

func TestNewRetroSpaceKey(t *testing.T) {
	spaceID := coretypes.SpaceID("s1")
	key := NewRetroSpaceKey(spaceID)
	if key.ID != const4retrospectus.ExtensionID {
		t.Errorf("key.ID = %v, want %v", key.ID, const4retrospectus.ExtensionID)
	}
}

func TestNewRetroSpaceEntry(t *testing.T) {
	spaceID := coretypes.SpaceID("s1")
	entry := NewRetroSpaceEntry(spaceID)
	if entry.ID != string(const4retrospectus.ExtensionID) {
		t.Errorf("entry.ID = %v, want %v", entry.ID, const4retrospectus.ExtensionID)
	}
	if entry.Data == nil {
		t.Error("entry.Data is nil")
	}
}
