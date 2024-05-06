package models4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/strongo/validation"
)

type CalendarHappeningBrief struct {
	HappeningBrief
	models4linkage.WithRelated
}

func (v *CalendarHappeningBrief) Validate() error {
	if err := v.HappeningBrief.Validate(); err != nil {
		return err
	}
	if err := v.WithRelated.Validate(); err != nil {
		return err
	}
	return nil
}

type CalendariumTeamDbo struct {
	RecurringHappenings map[string]*CalendarHappeningBrief `json:"recurringHappenings,omitempty" firestore:"recurringHappenings,omitempty"`
}

func (v *CalendariumTeamDbo) Validate() error {
	for i, rh := range v.RecurringHappenings {
		if err := rh.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("recurringHappenings", fmt.Errorf("invalid value at index %v: %w", i, err).Error())
		}
		if rh.Type != HappeningTypeRecurring {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("recurringHappenings[%v].type", i),
				fmt.Errorf("expected to have value 'recurring' got '%v'", rh.Type).Error())
		}
	}
	return nil
}

func (v *CalendariumTeamDbo) GetRecurringHappeningBrief(id string) *CalendarHappeningBrief {
	return v.RecurringHappenings[id]
}
