package dbo4calendarius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalendarDayDbo_GetAdjustment(t *testing.T) {
	slotAdj := &SlotAdjustment{}
	happeningAdj := &HappeningAdjustment{Slots: map[string]*SlotAdjustment{"s1": slotAdj}}
	dbo := CalendarDayDbo{
		HappeningAdjustments: map[string]*HappeningAdjustment{"h1": happeningAdj},
	}

	t.Run("happening_and_slot", func(t *testing.T) {
		ha, sa := dbo.GetAdjustment("h1", "s1")
		assert.Same(t, happeningAdj, ha)
		assert.Same(t, slotAdj, sa)
	})

	t.Run("happening_empty_slot_id", func(t *testing.T) {
		ha, sa := dbo.GetAdjustment("h1", "")
		assert.Same(t, happeningAdj, ha)
		assert.Nil(t, sa)
	})

	t.Run("missing_happening", func(t *testing.T) {
		ha, sa := dbo.GetAdjustment("missing", "s1")
		assert.Nil(t, ha)
		assert.Nil(t, sa)
	})

	t.Run("missing_slot", func(t *testing.T) {
		ha, sa := dbo.GetAdjustment("h1", "missing")
		assert.Same(t, happeningAdj, ha)
		assert.Nil(t, sa)
	})
}

func TestCalendarDayDbo_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dbo     CalendarDayDbo
		wantErr bool
	}{
		{
			name:    "valid",
			dbo:     CalendarDayDbo{HappeningIDs: []string{"h1"}},
			wantErr: false,
		},
		{
			name:    "empty_happening_ids",
			dbo:     CalendarDayDbo{},
			wantErr: true,
		},
		{
			name: "nil_adjustment",
			dbo: CalendarDayDbo{
				HappeningIDs:         []string{"h1"},
				HappeningAdjustments: map[string]*HappeningAdjustment{"h1": nil},
			},
			wantErr: true,
		},
		{
			name: "valid_adjustment",
			dbo: CalendarDayDbo{
				HappeningIDs:         []string{"h1"},
				HappeningAdjustments: map[string]*HappeningAdjustment{"h1": {}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dbo.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewCalendarDayKey(t *testing.T) {
	key := NewCalendarDayKey("space1", "2020-01-01")
	assert.NotNil(t, key)
	assert.Equal(t, "2020-01-01", key.ID)
}

func TestNewCalendarDayEntry(t *testing.T) {
	entry := NewCalendarDayEntry("space1", "2020-01-01")
	assert.Equal(t, "2020-01-01", entry.ID)
	assert.NotNil(t, entry.Key)
	assert.NotNil(t, entry.Data)
	assert.NotNil(t, entry.Record)
}

func TestNewCalendarDayEntry_Panics(t *testing.T) {
	assert.Panics(t, func() { NewCalendarDayEntry("", "2020-01-01") })
	assert.Panics(t, func() { NewCalendarDayEntry("space1", "not-a-date") })
}

func TestNewCalendarDayEntryWithDbo_PanicsOnNil(t *testing.T) {
	assert.Panics(t, func() { NewCalendarDayEntryWithDbo("space1", "2020-01-01", nil) })
}
