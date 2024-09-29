package dbo4calendarium

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

const ReasonMaxLen = 10000

const DaysCollection = "days"

// CalendarDayDbo is a record in "spaces/{teamID}/calendarium/days" collection with ContactID=YYYY-MM-DD
// It's needed to store adjustments of recurring happenings for a specific days like:
// - canceled slots
// - changed slot times
// - changed slot locations
// - changed slot participants
type CalendarDayDbo struct {
	HappeningIDs         []string                        `json:"happeningIDs,omitempty" firestore:"happeningIDs,omitempty"`
	HappeningAdjustments map[string]*HappeningAdjustment `json:"happeningAdjustments,omitempty" firestore:"happeningAdjustments,omitempty"`
}

// GetAdjustment returns adjustments for happening and slot for a specific day
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

// Validate returns error if not valid
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

// CalendarDayEntry is a helper type to work with CalendarDayDbo and it's key
type CalendarDayEntry = record.DataWithID[string, *CalendarDayDbo]

// NewCalendarDayKey returns key for a record in teams/{teamID}/calendarium/days collection with ContactID=YYYY-MM-DD
func NewCalendarDayKey(teamID, date string) *dal.Key {
	return dbo4spaceus.NewSpaceModuleItemKey(teamID, const4calendarium.ModuleID, DaysCollection, date)
}

// NewCalendarDayEntry creates a new instance of CalendarDayEntry
func NewCalendarDayEntry(teamID, date string) CalendarDayEntry {
	if teamID == "" {
		panic(errors.New("required parameter 'teamID' is empty string"))
	}
	if _, err := validate.DateString(date); err != nil {
		panic(err)
	}
	dto := new(CalendarDayDbo)
	return NewCalendarDayEntryWithDbo(teamID, date, dto)
}

// NewCalendarDayEntryWithDbo creates a new instance of CalendarDayEntry with provided CalendarDayDbo
func NewCalendarDayEntryWithDbo(teamID, date string, dbo *CalendarDayDbo) (calendarDay CalendarDayEntry) {
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
