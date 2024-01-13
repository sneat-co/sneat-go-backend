package facade4retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/models4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/models4retrospectus"
)

// MeetingRecordFactory factory
type MeetingRecordFactory struct {
}

// Collection "meetings"
func (MeetingRecordFactory) Collection() string {
	return "meetings"
}

// NewRecordData creates new api4meetingus record
func (MeetingRecordFactory) NewRecordData() models4meetingus.MeetingInstance {
	return &models4retrospectus.Retrospective{}
}
