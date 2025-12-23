package facade4scrumus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
)

func TestMeetingRecordFactory_NewRecord(t *testing.T) {
	recordFactory := MeetingRecordFactory{}
	record := recordFactory.NewRecordData()
	scrum := record.(*dbo4scrumus.Scrum)
	if scrum.Timer != nil {
		t.Error("Expects scrum.Timer to be nil")
		scrum.Timer = nil
	}
	meeting := record.BaseMeeting()
	if meeting.Timer != nil {
		t.Error("Expects api4meetingus.Timer to be nil")
		meeting.Timer = nil
	}
	meeting.Timer = new(dbo4meetingus.Timer)
	if scrum.Timer == nil {
		t.Error("Expects scrum.Timer to be NOT nil")
	}
}
