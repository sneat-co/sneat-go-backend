package dbo4meetingus

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func TestTimer_Validate(t *testing.T) {
	validBy := dbmodels.ByUser{UID: "u1"}
	tests := []struct {
		name    string
		timer   Timer
		wantErr bool
	}{
		{
			name:    "active",
			timer:   Timer{Status: TimerStatusActive, By: validBy, At: time.Now()},
			wantErr: false,
		},
		{
			name:    "stopped_valid",
			timer:   Timer{Status: TimerStatusStopped, By: validBy, ElapsedSeconds: 10},
			wantErr: false,
		},
		{
			name:    "stopped_without_elapsed",
			timer:   Timer{Status: TimerStatusStopped, By: validBy},
			wantErr: true,
		},
		{
			name:    "stopped_with_active_member",
			timer:   Timer{Status: TimerStatusStopped, By: validBy, ElapsedSeconds: 10, ActiveMemberID: "m1"},
			wantErr: true,
		},
		{
			name:    "paused_valid",
			timer:   Timer{Status: TimerStatusPaused, By: validBy, ElapsedSeconds: 5},
			wantErr: false,
		},
		{
			name:    "paused_without_elapsed",
			timer:   Timer{Status: TimerStatusPaused, By: validBy},
			wantErr: true,
		},
		{
			name:    "empty_status",
			timer:   Timer{Status: "", By: validBy},
			wantErr: true,
		},
		{
			name:    "unknown_status",
			timer:   Timer{Status: "bogus", By: validBy},
			wantErr: true,
		},
		{
			name:    "invalid_by",
			timer:   Timer{Status: TimerStatusActive, By: dbmodels.ByUser{}},
			wantErr: true,
		},
		{
			name:    "active_member_with_spaces",
			timer:   Timer{Status: TimerStatusActive, By: validBy, ActiveMemberID: " m1 "},
			wantErr: true,
		},
		{
			name:    "non_positive_seconds_by_member",
			timer:   Timer{Status: TimerStatusActive, By: validBy, SecondsByMember: map[string]int{"m1": 0}},
			wantErr: true,
		},
		{
			name:    "positive_seconds_by_member",
			timer:   Timer{Status: TimerStatusActive, By: validBy, SecondsByMember: map[string]int{"m1": 3}},
			wantErr: false,
		},
		{
			name:    "non_positive_seconds_by_topic",
			timer:   Timer{Status: TimerStatusActive, By: validBy, SecondsByTopic: map[string]int{"t1": -1}},
			wantErr: true,
		},
		{
			name:    "positive_seconds_by_topic",
			timer:   Timer{Status: TimerStatusActive, By: validBy, SecondsByTopic: map[string]int{"t1": 7}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.timer.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
