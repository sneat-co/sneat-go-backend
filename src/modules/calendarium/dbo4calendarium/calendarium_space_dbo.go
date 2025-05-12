package dbo4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/strongo/validation"
)

type CalendarHappeningBrief struct {
	HappeningBase
	dbo4linkage.WithRelated
}

func (v *CalendarHappeningBrief) Validate() error {
	if err := v.HappeningBase.Validate(); err != nil {
		return err
	}
	if err := v.WithRelated.Validate(); err != nil {
		return err
	}
	return nil
}

const RecurringHappeningsField = "recurringHappenings"
const SlotsField = "slots"

type CalendariumSpaceDbo struct {
	RecurringHappenings map[string]*CalendarHappeningBrief `json:"recurringHappenings,omitempty" firestore:"recurringHappenings,omitempty"`
}

func (v *CalendariumSpaceDbo) Validate() error {
	for i, rh := range v.RecurringHappenings {
		if err := rh.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(RecurringHappeningsField, fmt.Errorf("invalid value at index %v: %w", i, err).Error())
		}
		if rh.Type != HappeningTypeRecurring {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("recurringHappenings[%v].type", i),
				fmt.Errorf("expected to have value 'recurring' got '%v'", rh.Type).Error())
		}
	}
	return nil
}

func (v *CalendariumSpaceDbo) GetRecurringHappeningBrief(id string) *CalendarHappeningBrief {
	return v.RecurringHappenings[id]
}
