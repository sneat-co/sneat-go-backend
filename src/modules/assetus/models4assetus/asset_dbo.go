package models4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
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
	WithAssetExtraField
	briefs4assetus.WithAssetusTeamBriefs[*briefs4assetus.AssetBrief]
	with.TagsField
	briefs4contactus.WithMultiTeamContactIDs
	dbmodels.WithCustomFields
	AssetDates
}

func (v *AssetBaseDbo) Validate() error {
	if err := v.AssetBrief.Validate(); err != nil {
		return err
	}
	if err := v.WithAssetusTeamBriefs.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	if err := v.WithMultiTeamContactIDs.Validate(); err != nil {
		return err
	}
	if err := v.WithCustomFields.Validate(); err != nil {
		return err
	}
	if err := v.AssetDates.Validate(); err != nil {
		return err
	}
	if extra, err := v.GetExtra(); err != nil {
		return err
	} else if extra2, ok := extra.(WithAssetValidator); ok {
		if err := extra2.ValidateWithAsset(v); err != nil {
			return validation.NewErrBadRecordFieldValue("extra", err.Error())
		}
	} else if err := extra.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("extra", err.Error())
	}
	return nil
}

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
	if v.IsRequest {
		return validation.NewErrBadRecordFieldValue("isRequest", "should be false")
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
