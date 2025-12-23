package briefs4assetus

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sneat-co/sneat-core-modules/core/extra"

	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-core/geo"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

type AssetExtra interface {
	extra.Data
	ValidateWithAssetBrief(assetBrief AssetBrief) error
}

// AssetBrief keeps main props of an asset
type AssetBrief struct {
	Title      string                        `json:"title,omitempty" firestore:"title,omitempty"` // Should be required if the make, model & reg number are not provided
	Status     const4assetus.AssetStatus     `json:"status" firestore:"status"`                   // required field
	Category   const4assetus.AssetCategory   `json:"category" firestore:"category"`               // required field
	Type       const4assetus.AssetType       `json:"type" firestore:"type"`                       // required field
	Possession const4assetus.AssetPossession `json:"possession" firestore:"possession"`           // required field
	CountryID  geo.CountryAlpha2             `json:"countryID"  firestore:"countryID"`            // intentionally not omitempty so can be used in queries
	extra.WithExtraField
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
	if err := const4assetus.ValidateAssetType(v.Category, v.Type); err != nil {
		return err
	}

	if err := const4assetus.ValidateAssetPossession(v.Possession, true); err != nil {
		return err
	}
	if extraData, err := v.GetExtraData(); err != nil {
		return nil
	} else if assetExtra, ok := extraData.(AssetExtra); ok {
		if err = assetExtra.ValidateWithAssetBrief(*v); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("type %T does not implement AssetExtra interface", extraData)
	}
	if err := v.WithExtraField.Validate(); err != nil {
		return err
	}
	return nil
}
