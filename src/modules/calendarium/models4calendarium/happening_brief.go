package models4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// HappeningBrief hold data that stored both in entity record and in a brief.
type HappeningBrief struct {
	Type     HappeningType    `json:"type" firestore:"type"`
	Status   string           `json:"status" firestore:"status"`
	Canceled *Canceled        `json:"canceled,omitempty" firestore:"canceled,omitempty"`
	Title    string           `json:"title" firestore:"title"`
	Levels   []string         `json:"levels,omitempty" firestore:"levels,omitempty"`
	Slots    []*HappeningSlot `json:"slots,omitempty" firestore:"slots,omitempty"`

	//WithParticipants // TODO: replace with models4linkage.WithRelated

	// HappeningAssets keeps briefs for assets related to the happening.
	// Map key is expected to be valid dbmodels.TeamItemID to support contacts from multiple teams.
	HappeningAssets map[string]*HappeningAsset `json:"places,omitempty" firestore:"places,omitempty"`
}

func (v HappeningBrief) GetSlot(id string) (i int, slot *HappeningSlot) {
	for i, slot = range v.Slots {
		if slot.ID == id {
			return
		}
	}
	return -1, nil
}

// Validate returns error if not valid
func (v HappeningBrief) Validate() error {
	switch v.Type {
	case HappeningTypeSingle, HappeningTypeRecurring:
		break
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	default:
		return validation.NewErrBadRecordFieldValue("type", "unknown value: "+v.Type)
	}
	if v.Status == "" {
		return validation.NewErrRecordIsMissingRequiredField("status")
	}
	if !IsKnownHappeningStatus(v.Status) {
		return validation.NewErrBadRecordFieldValue("status", fmt.Sprintf("unknown value: '%v'", v.Status))
	}
	if v.Status == HappeningStatusCanceled && v.Canceled == nil {
		return validation.NewErrRecordIsMissingRequiredField("canceled")
	}
	if v.Canceled != nil && v.Status != HappeningStatusCanceled {
		return validation.NewErrBadRecordFieldValue("canceled", "should be populated only for canceled happenings, current status="+v.Status)
	}

	if err := dbmodels.ValidateTitle(v.Title); err != nil {
		return err
	}
	if len(v.Slots) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("slots")
	}

	for i, slot := range v.Slots {
		if err := slot.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("slots[%v]", i), err.Error())
		}
		for j, s := range v.Slots {
			if i != j && s.ID == slot.ID {
				return validation.NewErrBadRecordFieldValue("slots", fmt.Sprintf("at least 2 slots have same ContactID at indexes: %v & %v", i, j))
			}
			// TODO: Add more validations?
		}
	}

	//if err := v.WithParticipants.Validate(); err != nil {
	//	return err
	//}
	if err := validateHappeningAssetBriefs(v.HappeningAssets); err != nil {
		return err
	}

	return nil
}

func validateHappeningAssetBriefs(assets map[string]*HappeningAsset) error {
	for assetID, assetBrief := range assets {
		if assetID == "" {
			return validation.NewErrBadRecordFieldValue("happeningAssets", "assetID is empty")
		}
		field := func() string {
			return fmt.Sprintf("happeningAssets[%s]", assetID)
		}
		if err := dbmodels.TeamItemID(assetID).Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(field(), err.Error())
		}
		if err := assetBrief.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(field(), err.Error())
		}
	}
	return nil
}
