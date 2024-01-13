package models4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/strongoapp/with"
)

const TeamAssetsCollection = "assets"

type AssetSpecificData interface {
	Validate() error
}

type AssetMain interface {
	Validate() error
	AssetMainData() *AssetMainDto
	SpecificData() AssetSpecificData
	SetSpecificData(AssetSpecificData)
}

type AssetCreationData interface {
	Validate() error
	AssetMainData() *AssetMainDto
	SpecificData() AssetSpecificData
}

type AssetDto struct {
	with.TagsField
}

// AssetDbData defines mandatory fields & methods on an asset record
type AssetDbData interface {
	AssetMain
	AssetExtraData() *AssetExtraDto
}

// AssetMainDto was intended to be used in both AssetBaseDto and request to create an asset,
// but it was not a good idea as not clear how to manage module specific fields
type AssetMainDto struct {
	briefs4assetus.AssetBrief
	briefs4assetus.WithAssetusTeamBriefs[*briefs4assetus.AssetBrief]
	with.TagsField
	briefs4contactus.WithMultiTeamContactIDs
	dbmodels.WithCustomFields
	AssetDates
}

func (v *AssetMainDto) AssetMainData() *AssetMainDto {
	return v
}

func (v *AssetMainDto) Validate() error {
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
	return nil
}

// AssetExtraDto defines extra fields on an asset record that are not passed in create asset request
type AssetExtraDto struct {
	dbmodels.WithModified
	dbmodels.WithUserIDs
	dbmodels.WithTeamIDs
}

func (v *AssetExtraDto) AssetExtraData() *AssetExtraDto {
	return v
}

func (v *AssetExtraDto) Validate() error {
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
