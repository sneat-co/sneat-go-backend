package dbo4assetus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/core/extra"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
)

const TeamAssetsCollection = "assets"

type AssetCreationData interface {
	AssetBaseDboExtension
}

type WithAssetValidator interface {
	ValidateWithAsset(asset *AssetBaseDbo) error
}

type AssetBaseDboExtension interface {
	Validate() error
	GetAssetBaseDbo() *AssetBaseDbo
}

// AssetBaseDbo is used in both AssetDbo and a request to create an asset,
type AssetBaseDbo struct {
	briefs4assetus.AssetBrief
	WithAssetTeams
	with.TagsField
	//briefs4contactus.WithMultiTeamContactIDs
	dbmodels.WithCustomFields
	AssetDates
}

func (v *AssetBaseDbo) GetAssetBrief() (assetBrief briefs4assetus.AssetBrief, err error) {
	var extraData extra.Data
	if extraData, err = v.GetExtraData(); err != nil {
		return
	}
	assetBrief = v.AssetBrief
	extraData = extraData.GetBrief()
	if err = assetBrief.SetExtra(v.ExtraType, extraData); err != nil {
		return
	}
	return
}

func (v *AssetBaseDbo) Validate() error {
	if err := v.AssetBrief.Validate(); err != nil {
		return err
	}
	if err := v.WithAssetTeams.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	//if err := v.WithMultiTeamContactIDs.Validate(); err != nil {
	//	return err
	//}
	if err := v.WithCustomFields.Validate(); err != nil {
		return err
	}
	if err := v.AssetDates.Validate(); err != nil {
		return err
	}
	if extraData, err := v.GetExtraData(); err != nil {
		return err
	} else if extra2, ok := extraData.(WithAssetValidator); ok {
		if err := extra2.ValidateWithAsset(v); err != nil {
			return validation.NewErrBadRecordFieldValue("extraData", err.Error())
		}
	} else if err := extraData.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("extraData", err.Error())
	}
	return nil
}

type AssetEntry = record.DataWithID[string, *AssetDbo]

// AssetDbo defines fields on an asset record that are not passed in create asset request
type AssetDbo struct {
	AssetBaseDbo
	dbmodels.WithModified
	dbmodels.WithUserIDs
	dbmodels.WithTeamIDs
}

func (v *AssetDbo) AssetExtraData() *AssetDbo {
	return v
}

func (v *AssetDbo) Validate() error {
	if err := v.AssetBaseDbo.Validate(); err != nil {
		return err
	}
	if err := v.WithModified.Validate(); err != nil {
		return err
	}
	if err := v.WithUserIDs.Validate(); err != nil {
		return err
	}
	if err := v.WithTeamIDs.Validate(); err != nil {
		return err
	}
	return nil
}

// AssetDates defines dates of an asset - TODO: consider refactoring to custom fields?
type AssetDates struct {
	DateOfBuild       string `json:"dateOfBuild,omitempty" firestore:"dateOfBuild,omitempty"`
	DateOfPurchase    string `json:"dateOfPurchase,omitempty" firestore:"dateOfPurchase,omitempty"`
	DateInsuredTill   string `json:"dateInsuredTill,omitempty" firestore:"dateInsuredTill,omitempty"`
	DateCertifiedTill string `json:"dateCertifiedTill,omitempty" firestore:"dateCertifiedTill,omitempty"`
}

// Validate returns error if not valid
func (v *AssetDates) Validate() error {
	if v.DateOfBuild != "" {
		if _, err := validate.DateString(v.DateOfBuild); err != nil {
			return err
		}
	}
	if v.DateOfPurchase != "" {
		if _, err := validate.DateString(v.DateOfPurchase); err != nil {
			return err
		}
	}
	if v.DateInsuredTill != "" {
		if _, err := validate.DateString(v.DateInsuredTill); err != nil {
			return err
		}
	}
	if v.DateCertifiedTill != "" {
		if _, err := validate.DateString(v.DateCertifiedTill); err != nil {
			return err
		}
	}
	return nil
}
