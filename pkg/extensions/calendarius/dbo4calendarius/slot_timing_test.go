package dbo4calendarius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlotStatus_IsValid(t *testing.T) {
	tests := []struct {
		status SlotStatus
		want   bool
	}{
		{SlotStatus(SlotStatusActive), true},
		{SlotStatus(SlotStatusCanceled), true},
		{SlotStatus("unknown"), false},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestTiming_Validate(t *testing.T) {
	validStart := DateTime{Date: "2020-01-01", Time: "10:00"}
	validEnd := DateTime{Time: "11:00"}
	tests := []struct {
		name    string
		timing  Timing
		wantErr bool
	}{
		{"valid", Timing{Start: validStart, End: validEnd}, false},
		{"valid_with_duration", Timing{Start: validStart, End: validEnd, DurationInMinutes: 60}, false},
		{"invalid_start", Timing{Start: DateTime{}, End: validEnd}, true},
		{"invalid_end", Timing{Start: validStart, End: DateTime{Time: "bad"}}, true},
		{"negative_duration", Timing{Start: validStart, End: validEnd, DurationInMinutes: -1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.timing.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHappeningSlotTiming_IsEmpty(t *testing.T) {
	full := HappeningSlotTiming{
		Timing: Timing{
			Start: DateTime{Date: "2020-01-01", Time: "10:00"},
			End:   DateTime{Time: "11:00"},
		},
	}
	assert.False(t, full.IsEmpty(), "fully populated once-slot is not empty")

	noStartDate := full
	noStartDate.Start.Date = ""
	assert.True(t, noStartDate.IsEmpty())

	withWeekdays := full
	withWeekdays.Weekdays = []WeekdayCode{Monday2}
	assert.True(t, withWeekdays.IsEmpty())
}

func TestHappeningSlotTiming_Validate(t *testing.T) {
	onceTiming := Timing{
		Start: DateTime{Date: "2020-01-01", Time: "10:00"},
		End:   DateTime{Time: "11:00"},
	}
	weeklyTiming := Timing{
		Start: DateTime{Date: "2020-01-01", Time: "10:00"},
		End:   DateTime{Time: "11:00"},
	}
	tests := []struct {
		name    string
		v       HappeningSlotTiming
		wantErr bool
	}{
		{
			name:    "valid_once",
			v:       HappeningSlotTiming{Timing: onceTiming, Repeats: RepeatPeriodOnce},
			wantErr: false,
		},
		{
			name:    "valid_weekly",
			v:       HappeningSlotTiming{Timing: weeklyTiming, Repeats: RepeatPeriodWeekly, Weekdays: []WeekdayCode{Monday2}},
			wantErr: false,
		},
		{
			name:    "valid_daily",
			v:       HappeningSlotTiming{Timing: weeklyTiming, Repeats: RepeatPeriodDaily},
			wantErr: false,
		},
		{
			name:    "valid_monthly_with_weeks",
			v:       HappeningSlotTiming{Timing: weeklyTiming, Repeats: RepeatPeriodMonthly, Weeks: []int{1, 3}},
			wantErr: false,
		},
		{
			name:    "invalid_timing",
			v:       HappeningSlotTiming{Timing: Timing{Start: DateTime{}}, Repeats: RepeatPeriodOnce},
			wantErr: true,
		},
		{
			name:    "missing_repeats",
			v:       HappeningSlotTiming{Timing: onceTiming, Repeats: ""},
			wantErr: true,
		},
		{
			name:    "weekly_without_weekdays",
			v:       HappeningSlotTiming{Timing: weeklyTiming, Repeats: RepeatPeriodWeekly},
			wantErr: true,
		},
		{
			name:    "daily_with_weekdays",
			v:       HappeningSlotTiming{Timing: weeklyTiming, Repeats: RepeatPeriodDaily, Weekdays: []WeekdayCode{Monday2}},
			wantErr: true,
		},
		{
			name:    "unknown_repeats",
			v:       HappeningSlotTiming{Timing: weeklyTiming, Repeats: "fortnightly"},
			wantErr: true,
		},
		{
			name:    "once_missing_end_time",
			v:       HappeningSlotTiming{Timing: Timing{Start: DateTime{Date: "2020-01-01", Time: "10:00"}}, Repeats: RepeatPeriodOnce},
			wantErr: true,
		},
		{
			name:    "weekly_empty_weekday",
			v:       HappeningSlotTiming{Timing: weeklyTiming, Repeats: RepeatPeriodWeekly, Weekdays: []WeekdayCode{""}},
			wantErr: true,
		},
		{
			name:    "week_out_of_range",
			v:       HappeningSlotTiming{Timing: weeklyTiming, Repeats: RepeatPeriodMonthly, Weeks: []int{6}},
			wantErr: true,
		},
		{
			name:    "duplicate_week",
			v:       HappeningSlotTiming{Timing: weeklyTiming, Repeats: RepeatPeriodMonthly, Weeks: []int{1, 1}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
