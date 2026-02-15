package facade4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus/dbo4scrumus"
)

// MeetingRecordFactory factory
type MeetingRecordFactory struct {
}

// Collection "api4meetingus"
func (MeetingRecordFactory) Collection() string {
	return "api4meetingus"
}

// NewRecordData create new api4meetingus record
func (MeetingRecordFactory) NewRecordData() dbo4meetingus.MeetingInstance {
	return &dbo4scrumus.Scrum{}
}
