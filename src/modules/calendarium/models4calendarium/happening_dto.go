package models4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

// HappeningDbo DTO
type HappeningDbo struct {
	HappeningBrief
	with.CreatedFields
	with.TagsField
	dbmodels.WithUserIDs
	with.DatesFields
	models4linkage.WithRelatedAndIDs
	//dbmodels.WithTeamDates
	briefs4contactus.WithMultiTeamContacts[*briefs4contactus.ContactBrief]
}

// Validate returns error if not valid
func (v *HappeningDbo) Validate() error {
	if err := v.HappeningBrief.Validate(); err != nil {
		return err
	}
	if err := v.WithUserIDs.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	if err := v.DatesFields.Validate(); err != nil {
		return err
	}
	//if err := v.WithTeamDates.Validate(); err != nil {
	//	return err
	//}
	//if len(v.TeamIDs) == 0 {
	//	return validation.NewErrRecordIsMissingRequiredField("teamIDs")
	//}
	for i, level := range v.Levels {
		if l := strings.TrimSpace(level); l == "" {
			return validation.NewErrRecordIsMissingRequiredField(
				fmt.Sprintf("levels[%v]", i),
			)
		} else if l != level {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("levels[%v]", i),
				fmt.Sprintf("whitespaces at beginning or end: [%v]", level),
			)
		}
	}
	if err := v.WithMultiTeamContactIDs.Validate(); err != nil {
		return err
	}
	switch v.Type {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	case HappeningTypeSingle:
		if count := len(v.Slots); count > 1 {
			return validation.NewErrBadRecordFieldValue("slots", fmt.Sprintf("single time happening should have only single 'once' slot, got: %v", count))
		}
		if len(v.Dates) == 0 {
			return validation.NewErrRecordIsMissingRequiredField("dates")
		}
		//if len(v.TeamDates) == 0 {
		//	return validation.NewErrRecordIsMissingRequiredField("teamDates")
		//}
	case HappeningTypeRecurring:
		if len(v.Dates) > 0 {
			return validation.NewErrBadRequestFieldValue("dates", "should be empty for 'recurring' happening")
		}
	default:
		return validation.NewErrBadRecordFieldValue("type", "unknown value: "+v.Type)
	}

	if err := v.WithMultiTeamContacts.Validate(); err != nil {
		return err
	}
	return nil
}
