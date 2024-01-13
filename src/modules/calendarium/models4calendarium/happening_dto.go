package models4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

// HappeningDto DTO
type HappeningDto struct {
	HappeningBrief
	with.CreatedFields
	with.TagsField
	dbmodels.WithUserIDs
	with.DatesFields
	models4linkage.WithRelatedAndIDs
	//dbmodels.WithTeamDates
	briefs4contactus.WithMultiTeamContacts[*briefs4contactus.ContactBrief]
	AssetIDs []string `json:"assetIDs,omitempty" firestore:"assetIDs,omitempty"` // TODO: migrated to WithRelatedAndIDs
}

// Validate returns error if not valid
func (v *HappeningDto) Validate() error {
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
	if err := validateHappeningAssets(v.AssetIDs, v.HappeningAssets); err != nil {
		return err
	}
	//if v.Role == HappeningTypeRecurring && v.Status == HappeningStatusCanceled {
	//	for _, slot := range v.Slots {
	//		if slot.Status != SlotStatusCanceled {
	//
	//		}
	//	}
	//}
	return nil
}

func validateHappeningAssets(assetIDs []string, assets map[string]*HappeningAsset) error {
	if err := validateHappeningAssetIDs(assetIDs, assets); err != nil {
		return err
	}
	if err := validateHappeningAssetBriefs(assets); err != nil {
		return err
	}
	for assetID := range assets {
		if !slice.Contains(assetIDs, assetID) {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("happeningAssets[%s]", assetID),
				"asset ID is missing from assetIDs")
		}
	}
	return nil
}

func validateHappeningAssetIDs(assetIDs []string, assets map[string]*HappeningAsset) error {
	if len(assetIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("assetIDs")
	}
	if assetIDs[0] != "*" {
		return validation.NewErrBadRecordFieldValue("assetIDs[0]", "should be '*'")
	}
	for i, assetID := range assetIDs[1:] {
		if assetID == "" {
			return validation.NewErrBadRecordFieldValue("assetIDs", "assetID is empty")
		}
		field := func() string {
			return fmt.Sprintf("assetIDs[%d]", i)
		}
		if slice.Contains(assetIDs[:i], assetID) {
			return validation.NewErrBadRecordFieldValue(field(), "assetID is empty")
		}
		if err := dbmodels.TeamItemID(assetID).Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(field(), err.Error())
		}
		if _, ok := assets[assetID]; !ok {
			return validation.NewErrBadRecordFieldValue(field(), "asset brief is missing")
		}
	}
	return nil
}
