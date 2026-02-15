package facade4retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/retrospectus/dbo4retrospectus"
)

// MeetingRecordFactory factory
type MeetingRecordFactory struct {
}

// Collection "meetings"
func (MeetingRecordFactory) Collection() string {
	return "meetings"
}

// NewRecordData creates new api4meetingus record
func (MeetingRecordFactory) NewRecordData() dbo4meetingus.MeetingInstance {
	return &dbo4retrospectus.Retrospective{}
}
