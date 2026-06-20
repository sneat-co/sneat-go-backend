package dbo4scrumus

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/dbo4meetingus"
	"github.com/stretchr/testify/assert"
)

func TestScrum_BaseMeeting(t *testing.T) {
	v := &Scrum{}
	assert.Same(t, &v.Meeting, v.BaseMeeting())
}

func TestScrum_Validate_Timer(t *testing.T) {
	t.Run("invalid_timer", func(t *testing.T) {
		v := &Scrum{
			Meeting: dbo4meetingus.Meeting{
				Timer: &dbo4meetingus.Timer{Status: "bogus-status"},
			},
		}
		assert.Error(t, v.Validate())
	})
}

func TestScrum_Validate_ValidStatus(t *testing.T) {
	v := &Scrum{
		Statuses: ScrumStatusByMember{
			"m1": &MemberStatus{Member: ScrumMember{ID: "m1", Title: "Member 1"}},
		},
	}
	assert.NoError(t, v.Validate())
}

func TestScrumSettings_Validate(t *testing.T) {
	t.Run("nil_duration", func(t *testing.T) {
		v := &ScrumSettings{}
		assert.NoError(t, v.Validate())
	})
	t.Run("with_duration", func(t *testing.T) {
		v := &ScrumSettings{Duration: &dbo4spaceus.MeetingDurationSettings{}}
		assert.NoError(t, v.Validate())
	})
}

func TestScrumSpaceDto_Validate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := ScrumSpaceDto{}
		assert.NoError(t, v.Validate())
	})
	t.Run("with_settings", func(t *testing.T) {
		v := ScrumSpaceDto{ScrumSettings: &ScrumSettings{}}
		assert.NoError(t, v.Validate())
	})
	t.Run("invalid_active", func(t *testing.T) {
		v := ScrumSpaceDto{Active: &dbo4spaceus.SpaceMeetingInfo{}} // missing id/stage
		assert.Error(t, v.Validate())
	})
	t.Run("invalid_last", func(t *testing.T) {
		v := ScrumSpaceDto{Last: &dbo4spaceus.SpaceMeetingInfo{}} // missing id/stage
		assert.Error(t, v.Validate())
	})
	t.Run("valid_active_and_last", func(t *testing.T) {
		v := ScrumSpaceDto{
			Active: &dbo4spaceus.SpaceMeetingInfo{ID: "a", Stage: "active"},
			Last:   &dbo4spaceus.SpaceMeetingInfo{ID: "l", Stage: "ended"},
		}
		assert.NoError(t, v.Validate())
	})
}
