package dbo4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sharing/dbo4sharing"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

// HappeningDbo DTO
type HappeningDbo struct {
	HappeningBase
	with.CreatedFields
	with.TagsField
	dbmodels.WithUserIDs
	with.DatesFields
	dbo4linkage.WithRelatedAndIDs

	// TODO: Implement sharing!
	SharedTo *dbo4sharing.To `json:"sharedTo,omitempty" firestore:"sharedTo,omitempty"`

	Description string `json:"description,omitempty" firestore:"description,omitempty"`

	Adjustments HappeningAdjustment `json:"adjustments,omitempty" firestore:"adjustments,omitempty"`
	//dbmodels.WithSpaceDates
	//briefs4contactus.WithMultiSpaceContacts[*briefs4contactus.ContactBrief]
}

// Validate returns error if not valid
func (v *HappeningDbo) Validate() error {
	if err := v.HappeningBase.Validate(); err != nil {
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
	if err := v.Adjustments.Validate(); err != nil {
		return err
	}
	if v.SharedTo != nil {
		if err := v.SharedTo.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("sharedTo", err.Error())
		}
	}
	//if err := v.WithSpaceDates.Validate(); err != nil {
	//	return err
	//}
	//if len(v.SpaceIDs) == 0 {
	//	return validation.NewErrRecordIsMissingRequiredField("spaceIDs")
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
	switch v.Type {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	case HappeningTypeSingle:
		if count := len(v.Slots); count > 1 {
			return validation.NewErrBadRecordFieldValue(SlotsField, fmt.Sprintf("single time happening should have only single 'once' slot, got: %v", count))
		}
		if len(v.Dates) == 0 {
			return validation.NewErrRecordIsMissingRequiredField("dates")
		}
		//if len(v.TeamDates) == 0 {
		//	return validation.NewErrRecordIsMissingRequiredField("spaceDates")
		//}
	case HappeningTypeRecurring:
		if len(v.Dates) > 0 {
			return validation.NewErrBadRequestFieldValue("dates", "should be empty for 'recurring' happening")
		}
	default:
		return validation.NewErrBadRecordFieldValue("type", "unknown value: "+v.Type)
	}

	if err := v.WithRelatedIDs.Validate(); err != nil {
		return err
	}
	if err := dbo4linkage.ValidateRelatedAndRelatedIDs(v.WithRelated, v.RelatedIDs); err != nil {
		return err
	}

	//if err := v.WithMultiSpaceContacts.Validate(); err != nil {
	//	return err
	//}
	return nil
}
