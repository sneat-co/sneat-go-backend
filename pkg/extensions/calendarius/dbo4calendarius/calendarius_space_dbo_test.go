package dbo4calendarius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func validRecurringBrief() *CalendarHappeningBrief {
	return &CalendarHappeningBrief{
		HappeningBase: HappeningBase{
			Type:   HappeningTypeRecurring,
			Status: HappeningStatusActive,
			Title:  "Weekly class",
			Slots: map[string]*HappeningSlot{
				"s1": {HappeningSlotTiming: HappeningSlotTiming{
					Timing:   Timing{Start: DateTime{Date: "2020-01-01", Time: "10:00"}, End: DateTime{Time: "11:00"}},
					Repeats:  RepeatPeriodWeekly,
					Weekdays: []WeekdayCode{Monday2},
				}},
			},
		},
	}
}

func TestCalendarHappeningBrief_Validate(t *testing.T) {
	assert.NoError(t, validRecurringBrief().Validate())

	bad := validRecurringBrief()
	bad.Title = ""
	assert.Error(t, bad.Validate())
}

func TestCalendariusSpaceDbo_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dbo     *CalendariusSpaceDbo
		wantErr bool
	}{
		{"empty", &CalendariusSpaceDbo{}, false},
		{
			name:    "valid_recurring",
			dbo:     &CalendariusSpaceDbo{RecurringHappenings: map[string]*CalendarHappeningBrief{"h1": validRecurringBrief()}},
			wantErr: false,
		},
		{
			name: "invalid_brief",
			dbo: func() *CalendariusSpaceDbo {
				b := validRecurringBrief()
				b.Title = ""
				return &CalendariusSpaceDbo{RecurringHappenings: map[string]*CalendarHappeningBrief{"h1": b}}
			}(),
			wantErr: true,
		},
		{
			name: "wrong_type",
			dbo: func() *CalendariusSpaceDbo {
				b := validRecurringBrief()
				b.Type = HappeningTypeSingle
				return &CalendariusSpaceDbo{RecurringHappenings: map[string]*CalendarHappeningBrief{"h1": b}}
			}(),
			wantErr: true,
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

func TestCalendariusSpaceDbo_GetRecurringHappeningBrief(t *testing.T) {
	brief := validRecurringBrief()
	dbo := &CalendariusSpaceDbo{RecurringHappenings: map[string]*CalendarHappeningBrief{"h1": brief}}
	assert.Same(t, brief, dbo.GetRecurringHappeningBrief("h1"))
	assert.Nil(t, dbo.GetRecurringHappeningBrief("missing"))
}

func TestHappeningBrief_Validate(t *testing.T) {
	v := &HappeningBrief{HappeningBase: *validHappeningBase()}
	assert.NoError(t, v.Validate())

	bad := &HappeningBrief{HappeningBase: *validHappeningBase()}
	bad.Title = ""
	assert.Error(t, bad.Validate())
}
