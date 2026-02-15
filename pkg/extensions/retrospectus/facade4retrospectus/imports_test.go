package facade4retrospectus

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
	"testing"
)

func TestNewSpaceKey(t *testing.T) {
	spaceID := coretypes.SpaceID("s1")
	key := newSpaceKey(spaceID)
	if key.ID != string(spaceID) {
		t.Errorf("key.ID = %v, want %v", key.ID, spaceID)
	}
}

func TestMeetingRecordFactory(t *testing.T) {
	f := MeetingRecordFactory{}
	if f.Collection() != "meetings" {
		t.Errorf("Collection() = %v, want meetings", f.Collection())
	}
	if f.NewRecordData() == nil {
		t.Error("NewRecordData() returned nil")
	}
}
