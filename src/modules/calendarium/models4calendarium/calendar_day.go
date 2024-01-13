package models4calendarium

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strings"
)

// HappeningAdjustment at the moment supposed to be used only for recurring happenings
type HappeningAdjustment struct {
	HappeningID string        `json:"happeningID" firestore:"happeningID"`
	Slot        HappeningSlot `json:"slot" firestore:"slot"`
	Canceled    *Canceled     `json:"canceled,omitempty" firestore:"canceled,omitempty"`
}

const ReasonMaxLen = 10000

func (v *HappeningAdjustment) Validate() error {
	if v == nil {
		return errors.New("nil")
	}
	if strings.TrimSpace(v.HappeningID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("happeningID")
	}
	if err := v.Slot.Validate(); err != nil {
		return err
	}
	if v.Slot.Repeats == "recurring" {
		return validation.NewErrBadRecordFieldValue("slot.repeats", fmt.Sprintf("must be 'once', got '%v'", v.Slot.Repeats))
	}
	if v.Canceled != nil {
		if err := v.Canceled.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("canceled", err.Error())
		}
	}
	return nil
}

const CalendarDayCollection = "calendar_days"

type CalendarDayDto struct {
	dbmodels.WithTeamID
	Date                 string                 `json:"date" firestore:"date"`
	HappeningIDs         []string               `json:"happeningIDs" firestore:"happeningIDs"`
	HappeningAdjustments []*HappeningAdjustment `json:"happeningAdjustments" firestore:"happeningAdjustments"`
	//Happenings    []*HappeningBrief                 `json:"happenings" firestore:"happenings"`
}

func (v CalendarDayDto) GetAdjustment(happeningID, slotID string) (i int, adjustment *HappeningAdjustment) {
	for i, adjustment = range v.HappeningAdjustments {
		if adjustment.HappeningID == happeningID && adjustment.Slot.ID == slotID {
			return i, adjustment
		}
	}
	return -1, nil
}

func (v CalendarDayDto) Validate() error {
	if err := v.WithTeamID.Validate(); err != nil {
		return err
	}
	if v.Date == "" {
		return validation.NewErrRecordIsMissingRequiredField("date")
	}
	if _, err := validate.DateString(v.Date); err != nil {
		return validation.NewErrBadRecordFieldValue("date", err.Error())
	}
	if len(v.HappeningIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("happeningIDs")
	}
	for i, adjustment := range v.HappeningAdjustments {
		if err := adjustment.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("happeningAdjustments[%v]", i), err.Error())
		}
	}
	return nil
}

type CalendarDayEntry struct {
	record.WithID[string]
	Dto *CalendarDayDto
}

func (v CalendarDayEntry) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if v.Dto == nil {
		return validation.NewErrRecordIsMissingRequiredField("dto")
	}
	return v.Dto.Validate()
}

func NewCalendarDayID(teamID, date string) string {
	return teamID + ":" + date
}

func NewCalendarDayKey(teamID, date string) *dal.Key {
	id := NewCalendarDayID(teamID, date)
	return dal.NewKeyWithID(CalendarDayCollection, id)
}

func NewCalendarDayContext(teamID, date string) CalendarDayEntry {
	if teamID == "" {
		panic(errors.New("required parameter 'teamID' is empty string"))
	}
	if _, err := validate.DateString(date); err != nil {
		panic(err)
	}
	dto := new(CalendarDayDto)
	dto.TeamID = teamID
	dto.Date = date
	return NewCalendarDayContextWithDto(dto)
}

func NewCalendarDayContextWithDto(dto *CalendarDayDto) (calendarDay CalendarDayEntry) {
	if dto == nil {
		panic("dto is nil")
	}
	if dto.TeamID == "" {
		panic("dto.TeamID is empty string")
	}
	if dto.Date == "" {
		panic("dto.Date is empty string")
	}
	key := NewCalendarDayKey(dto.TeamID, dto.Date)
	calendarDay.ID = dto.Date
	calendarDay.FullID = NewCalendarDayID(dto.TeamID, dto.Date)
	calendarDay.Key = key
	calendarDay.Dto = dto
	calendarDay.Record = dal.NewRecordWithData(key, dto)
	return
}
