package facade4meetingus

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func TestToggleTimerResponse_Validate(t *testing.T) {
	validTimer := func() *dbo4meetingus.Timer {
		return &dbo4meetingus.Timer{
			Status: dbo4meetingus.TimerStatusActive,
			By:     dbmodels.ByUser{UID: "user1"},
			At:     time.Now(),
		}
	}
	tests := []struct {
		name    string
		timer   *dbo4meetingus.Timer
		wantErr bool
	}{
		{name: "valid_active", timer: validTimer(), wantErr: false},
		{name: "missing_status", timer: &dbo4meetingus.Timer{By: dbmodels.ByUser{UID: "user1"}, At: time.Now()}, wantErr: true},
		{name: "unknown_status", timer: &dbo4meetingus.Timer{Status: "bogus", By: dbmodels.ByUser{UID: "user1"}, At: time.Now()}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ToggleTimerResponse{Timer: tt.timer}
			err := resp.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
