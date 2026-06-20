package dbo4calendarius

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/with"
)

func validSingleHappeningDbo() *HappeningDbo {
	dbo := &HappeningDbo{
		HappeningBase: *validHappeningBase(),
	}
	dbo.UserIDs = []string{"u1"}
	dbo.Dates = []string{"2020-01-01"}
	dbo.RelatedIDs = []string{"-"}
	return dbo
}

func TestHappeningDbo_Validate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(v *HappeningDbo)
		wantErr bool
	}{
		{"valid_single", func(v *HappeningDbo) {}, false},
		{
			name:    "invalid_base",
			modify:  func(v *HappeningDbo) { v.Title = "" },
			wantErr: true,
		},
		{
			name:    "missing_user_ids",
			modify:  func(v *HappeningDbo) { v.WithUserIDs = dbmodels.WithUserIDs{} },
			wantErr: true,
		},
		{
			name:    "single_without_dates",
			modify:  func(v *HappeningDbo) { v.DatesFields = with.DatesFields{} },
			wantErr: true,
		},
		{
			name: "single_with_two_slots",
			modify: func(v *HappeningDbo) {
				v.Slots["s2"] = &HappeningSlot{HappeningSlotTiming: validSlotTiming()}
			},
			wantErr: true,
		},
		{
			name: "level_with_whitespace",
			modify: func(v *HappeningDbo) {
				v.Levels = []string{" lvl "}
			},
			wantErr: true,
		},
		{
			name: "empty_level",
			modify: func(v *HappeningDbo) {
				v.Levels = []string{"  "}
			},
			wantErr: true,
		},
		{
			name: "recurring_with_dates",
			modify: func(v *HappeningDbo) {
				v.Type = HappeningTypeRecurring
				v.Slots["s1"].Repeats = RepeatPeriodWeekly
				v.Slots["s1"].Weekdays = []WeekdayCode{Monday2}
			},
			wantErr: true,
		},
		{
			name: "valid_recurring",
			modify: func(v *HappeningDbo) {
				v.Type = HappeningTypeRecurring
				v.Slots["s1"].Repeats = RepeatPeriodWeekly
				v.Slots["s1"].Weekdays = []WeekdayCode{Monday2}
				v.DatesFields = with.DatesFields{}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validSingleHappeningDbo()
			tt.modify(v)
			err := v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
