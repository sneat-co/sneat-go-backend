package briefs4assetus

import (
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/extras4assetus"
	"strings"

	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-core/geo"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
)

// AssetBrief keeps main props of an asset
type AssetBrief struct {
	Title      string                        `json:"title" firestore:"title"`           // Should be required if the make, model & reg number are not provided
	Status     const4assetus.AssetStatus     `json:"status" firestore:"status"`         // required field
	Category   const4assetus.AssetCategory   `json:"category" firestore:"category"`     // required field
	Type       const4assetus.AssetType       `json:"type" firestore:"type"`             // required field
	Possession const4assetus.AssetPossession `json:"possession" firestore:"possession"` // required field
	CountryID  geo.CountryAlpha2             `json:"countryID"  firestore:"countryID"`  // intentionally not omitempty so can be used in queries
	extras4assetus.WithAssetExtraField
	dbmodels.WithOptionalRelatedAs
}

//func (v *AssetBrief) Equal(v2 *AssetBrief) bool {
//	return *v == *v2
//}

// Validate returns error if not valid
func (v *AssetBrief) Validate() error {
	if v == nil {
		return errors.New("can not be nil")
	}

	if err := v.WithOptionalRelatedAs.Validate(); err != nil {
		return err
	}
	if !const4assetus.IsValidAssetStatus(v.Status) {
		return validation.NewErrBadRecordFieldValue("status", fmt.Sprintf("unknown status: %s", v.Status))
	}
	if strings.TrimSpace(v.CountryID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("countryID")
	}
	checkType := func(types []string) error {
		switch v.Type {
		case "":
			return validation.NewErrRecordIsMissingRequiredField("type")
		default:
			if slice.Index(types, v.Type) < 0 {
				return validation.NewErrBadRecordFieldValue("type", fmt.Sprintf("unknown %s type: %s", v.Category, v.Type))
			}
		}
		return nil
	}
	switch v.Category {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("category")
	case const4assetus.AssetCategoryVehicle:
		//if strings.TrimSpace(v.Make) == "" {
		//	return validation.NewErrRecordIsMissingRequiredField("make")
		//}
		//if strings.TrimSpace(v.Model) == "" {
		//	return validation.NewErrRecordIsMissingRequiredField("model")
		//}
		//if err := checkType(const4assetus.AssetVehicleTypes); err != nil {
		//	return err
		//}
	case const4assetus.AssetCategoryDwelling:
		if err := checkType(const4assetus.AssetRealEstateTypes); err != nil {
			return err
		}
	case const4assetus.AssetCategorySportGear:
		if err := checkType(const4assetus.AssetSportGearTypes); err != nil {
			return err
		}
	case const4assetus.AssetCategoryDocument:
		if err := checkType(const4assetus.AssetDocumentTypes); err != nil {
			return err
		}
	default:
		return validation.NewErrBadRecordFieldValue("category", "unknown asset category: "+string(v.Category))
	}

	if err := const4assetus.ValidateAssetPossession(v.Possession, true); err != nil {
		return err
	}
	return nil
}
