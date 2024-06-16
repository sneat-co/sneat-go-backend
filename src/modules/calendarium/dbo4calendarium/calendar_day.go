package dbo4calendarium

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

const ReasonMaxLen = 10000

const DaysCollection = "days"

type CalendarDayDbo struct {
	HappeningIDs         []string                        `json:"happeningIDs,omitempty" firestore:"happeningIDs,omitempty"`
	HappeningAdjustments map[string]*HappeningAdjustment `json:"happeningAdjustments,omitempty" firestore:"happeningAdjustments,omitempty"`
}

func (v CalendarDayDbo) GetAdjustment(happeningID, slotID string) (
	happeningAdjustment *HappeningAdjustment,
	slotAdjustment *SlotAdjustment,
) {

	if happeningAdjustment = v.HappeningAdjustments[happeningID]; happeningAdjustment != nil && len(happeningAdjustment.Slots) > 0 {
		if slotID == "" {
			return happeningAdjustment, nil
		}
		return happeningAdjustment, happeningAdjustment.Slots[slotID]
	}
	return happeningAdjustment, nil
}

func (v CalendarDayDbo) Validate() error {
	if len(v.HappeningIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("happeningIDs")
	}
	for happeningID, adjustment := range v.HappeningAdjustments {
		if adjustment == nil {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("happeningAdjustments[%v]", happeningID))
		}
		if err := adjustment.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("happeningAdjustments[%v]", happeningID), err.Error())
		}
	}
	return nil
}

type CalendarDayEntry = record.DataWithID[string, *CalendarDayDbo]

func NewCalendarDayKey(teamID, date string) *dal.Key {
	return dal4teamus.NewTeamModuleItemKey(teamID, const4calendarium.ModuleID, DaysCollection, date)
}

func NewCalendarDayContext(teamID, date string) CalendarDayEntry {
	if teamID == "" {
		panic(errors.New("required parameter 'teamID' is empty string"))
	}
	if _, err := validate.DateString(date); err != nil {
		panic(err)
	}
	dto := new(CalendarDayDbo)
	return NewCalendarDayContextWithDbo(teamID, date, dto)
}

func NewCalendarDayContextWithDbo(teamID, date string, dbo *CalendarDayDbo) (calendarDay CalendarDayEntry) {
	if dbo == nil {
		panic("dbo is nil")
	}
	key := NewCalendarDayKey(teamID, date)
	calendarDay.ID = date
	calendarDay.Key = key
	calendarDay.Data = dbo
	calendarDay.Record = dal.NewRecordWithData(key, dbo)
	return
}
