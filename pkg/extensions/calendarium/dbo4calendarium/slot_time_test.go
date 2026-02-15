package dbo4calendarium

import "testing"

func TestDateTime_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dt      DateTime
		wantErr bool
	}{
		{"valid_date_only", DateTime{Date: "2020-01-01"}, false},
		{"valid_time_only", DateTime{Time: "12:00"}, false},
		{"valid_both", DateTime{Date: "2020-01-01", Time: "12:00"}, false},
		{"empty", DateTime{}, true},
		{"invalid_date", DateTime{Date: "invalid"}, true},
		{"invalid_time", DateTime{Time: "invalid"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dt.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DateTime.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsKnownRepeatPeriod(t *testing.T) {
	tests := []struct {
		period RepeatPeriod
		want   bool
	}{
		{RepeatPeriodOnce, true},
		{RepeatPeriodDaily, true},
		{"invalid", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(string(tt.period), func(t *testing.T) {
			if got := IsKnownRepeatPeriod(tt.period); got != tt.want {
				t.Errorf("IsKnownRepeatPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}
