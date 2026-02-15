package dbo4sportus

import (
	"testing"
	"time"
)

func TestSpotBrief_Validate(t *testing.T) {
	tests := []struct {
		name    string
		brief   SpotBrief
		wantErr bool
	}{
		{"valid", SpotBrief{Title: "Spot 1"}, false},
		{"empty_title", SpotBrief{Title: ""}, true},
		{"whitespace_title", SpotBrief{Title: "  "}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.brief.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SpotBrief.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSpotVisit_Validate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		visit   SpotVisit
		wantErr bool
	}{
		{"valid", SpotVisit{
			UserID:      "u1",
			SpotID:      "s1",
			CheckedInAt: now,
		}, false},
		{"missing_user", SpotVisit{
			SpotID:      "s1",
			CheckedInAt: now,
		}, true},
		{"missing_spot", SpotVisit{
			UserID:      "u1",
			CheckedInAt: now,
		}, true},
		{"invalid_dates", SpotVisit{
			UserID:       "u1",
			SpotID:       "s1",
			CheckedInAt:  now,
			CheckedOutAt: now.Add(-time.Hour),
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.visit.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SpotVisit.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
