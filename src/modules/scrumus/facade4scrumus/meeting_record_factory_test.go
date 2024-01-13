package facade4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/models4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/models4scrumus"
	"testing"
)

func TestMeetingRecordFactory_NewRecord(t *testing.T) {
	recordFactory := MeetingRecordFactory{}
	record := recordFactory.NewRecordData()
	scrum := record.(*models4scrumus.Scrum)
	if scrum.Timer != nil {
		t.Error("Expects scrum.Timer to be nil")
		scrum.Timer = nil
	}
	meeting := record.BaseMeeting()
	if meeting.Timer != nil {
		t.Error("Expects api4meetingus.Timer to be nil")
		meeting.Timer = nil
	}
	meeting.Timer = new(models4meetingus.Timer)
	if scrum.Timer == nil {
		t.Error("Expects scrum.Timer to be NOT nil")
	}
}
