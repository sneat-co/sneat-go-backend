package models4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strconv"
)

// Repeat defines repeat mode
type Repeat = string

const (
	// RepeatOnce = "once"
	RepeatOnce Repeat = "once"

	// RepeatDaily = "daily"
	RepeatDaily Repeat = "daily"

	// RepeatWeekly = "weekly"
	RepeatWeekly Repeat = "weekly"

	// RepeatMonthly = "monthly"
	RepeatMonthly Repeat = "monthly"

	// RepeatYearly = "yearly"
	RepeatYearly Repeat = "yearly"
)

// WeekdayCode defines weekday 2 chars code. Values: mo, tu, we, th, fr, sa, su
type WeekdayCode = string

const (
	// Monday2 = "mo"
	Monday2 WeekdayCode = "mo"

	// Tuesday2 = "tu"
	Tuesday2 WeekdayCode = "tu"

	// Wednesday2 = "we"
	Wednesday2 WeekdayCode = "we"

	// Thursday2 = "th"
	Thursday2 WeekdayCode = "th"

	// Friday2 = "fr"
	Friday2 WeekdayCode = "fr"

	// Saturday2 = "sa"
	Saturday2 WeekdayCode = "sa"

	// Sunday2 = "su"
	Sunday2 WeekdayCode = "su"
)

// DateTime DTO
type DateTime struct {
	Date string `json:"date,omitempty" firestore:"date,omitempty"`
	Time string `json:"time,omitempty" firestore:"time,omitempty"`
}

// Validate returns error if not valid
func (v DateTime) Validate() error {
	if v.Date != "" {
		if _, err := validate.DateString(v.Date); err != nil {
			return validation.NewErrBadRecordFieldValue("date", err.Error())
		}
	}
	if v.Time != "" {
		if err := validate.IsValidateTime(v.Time); err != nil {
			return validation.NewErrBadRecordFieldValue("time", err.Error())
		}
	}
	if v.Date == "" && v.Time == "" {
		return validation.NewErrBadRecordFieldValue("date|time", "either date or time or both fields should have value")
	}
	return nil
}

const (
	SlotStatusActive   = ""
	SlotStatusCanceled = "canceled"
)

func IsKnownSlotStatus(status string) bool {
	switch status {
	case SlotStatusActive, SlotStatusCanceled:
		return true
	}
	return false
}

type SlotStatus string

func (v SlotStatus) IsValid() bool {
	return IsKnownSlotStatus((string)(v))
}

// Timing DTO
type Timing struct {
	Start             DateTime `json:"start" firestore:"start"` // is required field
	End               DateTime `json:"end,omitempty" firestore:"end,omitempty"`
	DurationInMinutes int      `json:"durationInMinutes,omitempty" firestore:"durationInMinutes,omitempty"`
}

// Validate returns error if not valid
func (v Timing) Validate() error {
	if err := v.Start.Validate(); err != nil {
		return err
	}
	if err := v.End.Validate(); err != nil {
		return err
	}
	if v.DurationInMinutes < 0 {
		return validation.NewErrBadRecordFieldValue("durationInMinutes", "should be positive, got: "+strconv.Itoa(v.DurationInMinutes))
	}
	return nil
}

// HappeningSlotTiming DTO
type HappeningSlotTiming struct {
	Timing
	Repeats Repeat `json:"repeats" firestore:"repeats"`

	Weekdays []WeekdayCode `json:"weekdays,omitempty" firestore:"weekdays,omitempty"`

	// e.g. with [1,3]: repeats=monthly => every 1 & 3d week of every month.
	Weeks []int `json:"weeks,omitempty" firestore:"weeks,omitempty"`
}

// Validate returns error if not valid
func (v HappeningSlotTiming) Validate() error {
	if err := v.Timing.Validate(); err != nil {
		return err
	}
	switch v.Repeats {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("repeats")
	case RepeatWeekly:
		if len(v.Weekdays) == 0 {
			return validation.NewErrBadRecordFieldValue("weekdays", "for weekly recurring happenings weekdays also should be specified")
		}
	case RepeatOnce, RepeatDaily, RepeatMonthly, RepeatYearly:
		if len(v.Weekdays) > 0 {
			return validation.NewErrBadRecordFieldValue("weekdays", "can be specified only for weekly recurring happenings")
		}
	default:
		return validation.NewErrBadRecordFieldValue("repeats", "unknown value: "+v.Repeats)
	}
	if v.Repeats == RepeatOnce {
		if v.Start.Date == "" {
			return validation.NewErrRecordIsMissingRequiredField("slots[0].start.date")
		}
		if v.Start.Time == "" {
			return validation.NewErrRecordIsMissingRequiredField("slots[0].start.time")
		}
		if v.End.Time == "" {
			return validation.NewErrRecordIsMissingRequiredField("slots[0].end.time")
		}
	}
	for i, wd := range v.Weekdays {
		switch wd {
		case "":
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("weekdays[%v]", i))
		case Monday2, Tuesday2, Wednesday2, Thursday2, Friday2, Saturday2, Sunday2:
			break
		default:

		}
	}
	return nil
}
