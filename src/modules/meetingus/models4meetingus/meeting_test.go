package models4meetingus

import "testing"

func TestMeeting_HasUserID(t *testing.T) {
	meeting := Meeting{}
	t.Run("empty_record", func(t *testing.T) {
		if meeting.HasUserID("") {
			t.Error("expected false for empty UserID in empty record")
		}
	})
}

func TestMeeting_Validate(t *testing.T) {
	meeting := Meeting{}
	if err := meeting.Validate(); err != nil {
		t.Fatalf("no error expected for empty value, got: %v", err)
	}
}
