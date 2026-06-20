package dbo4calendarius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func validSlotTiming() HappeningSlotTiming {
	return HappeningSlotTiming{
		Timing: Timing{
			Start: DateTime{Date: "2020-01-01", Time: "10:00"},
			End:   DateTime{Time: "11:00"},
		},
		Repeats: RepeatPeriodOnce,
	}
}

func TestLocation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		loc     Location
		wantErr bool
	}{
		{"valid_physical", Location{Type: "physical", Title: "Office"}, false},
		{"valid_online", Location{Type: "online", Title: "Zoom"}, false},
		{"missing_type", Location{Title: "Office"}, true},
		{"unknown_type", Location{Type: "virtual", Title: "Office"}, true},
		{"missing_title", Location{Type: "physical"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.loc.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHappeningSlot_IsEmpty(t *testing.T) {
	var nilSlot *HappeningSlot
	assert.True(t, nilSlot.IsEmpty())

	empty := &HappeningSlot{}
	assert.True(t, empty.IsEmpty())

	withLocation := &HappeningSlot{Locations: []Location{{Type: "online", Title: "Zoom"}}}
	assert.False(t, withLocation.IsEmpty())
}

func TestHappeningSlot_Validate(t *testing.T) {
	tests := []struct {
		name    string
		slot    *HappeningSlot
		wantErr bool
	}{
		{
			name:    "valid",
			slot:    &HappeningSlot{HappeningSlotTiming: validSlotTiming()},
			wantErr: false,
		},
		{
			name: "valid_with_location",
			slot: &HappeningSlot{
				HappeningSlotTiming: validSlotTiming(),
				Locations:           []Location{{Type: "physical", Title: "Office"}},
			},
			wantErr: false,
		},
		{
			name: "invalid_location",
			slot: &HappeningSlot{
				HappeningSlotTiming: validSlotTiming(),
				Locations:           []Location{{Type: "bad", Title: "Office"}},
			},
			wantErr: true,
		},
		{
			name: "two_physical_locations",
			slot: &HappeningSlot{
				HappeningSlotTiming: validSlotTiming(),
				Locations: []Location{
					{Type: "physical", Title: "Office 1"},
					{Type: "physical", Title: "Office 2"},
				},
			},
			wantErr: true,
		},
		{
			name:    "invalid_timing",
			slot:    &HappeningSlot{HappeningSlotTiming: HappeningSlotTiming{Repeats: ""}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.slot.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
