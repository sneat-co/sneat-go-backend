package dal4scrumus

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
	"testing"
)

func TestNewScrumSpaceKey(t *testing.T) {
	spaceID := coretypes.SpaceID("s1")
	key := NewScrumSpaceKey(spaceID)
	if key.ID != ScrumusModuleID {
		t.Errorf("key.ID = %v, want %v", key.ID, ScrumusModuleID)
	}
}

func TestNewScrumSpaceEntry(t *testing.T) {
	spaceID := coretypes.SpaceID("s1")
	entry := NewScrumSpaceEntry(spaceID)
	if entry.ID != ScrumusModuleID {
		t.Errorf("entry.ID = %v, want %v", entry.ID, ScrumusModuleID)
	}
	if entry.Data == nil {
		t.Error("entry.Data is nil")
	}
}
