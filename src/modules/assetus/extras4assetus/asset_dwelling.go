package extras4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodels/extra"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"
	"github.com/strongo/validation"
)

func init() {
	RegisterAssetExtraFactory(AssetExtraTypeDwelling, func() briefs4assetus.AssetExtra {
		return new(AssetDwellingExtra)
	})
}

var _ extra.Data = (*AssetDwellingExtra)(nil)
var _ briefs4assetus.AssetExtra = (*AssetDwellingExtra)(nil)

type AssetDwellingExtra struct {
	//extra.BaseData
	Address   string `json:"address,omitempty" firestore:"address,omitempty"`
	RentPrice struct {
		Value    float64 `json:"value,omitempty" firestore:"value,omitempty"`
		Currency string  `json:"currency,omitempty" firestore:"currency,omitempty"`
	} `json:"rent_price,omitempty" firestore:"rent_price,omitempty"`
	NumberOfBedrooms int `json:"numberOfBedrooms,omitempty" firestore:"numberOfBedrooms,omitempty"`
	AreaSqM          int `json:"areaSqM,omitempty" firestore:"areaSqM,omitempty"`
}

func (v AssetDwellingExtra) ValidateWithAssetBrief(assetBrief briefs4assetus.AssetBrief) error {
	if err := v.Validate(); err != nil {
		return err
	}
	if assetBrief.Title == "" && v.Address == "" {
		return validation.NewValidationError("dwelling asset should have at least 1 of next fields: title, address")
	}
	return nil
}

func (v AssetDwellingExtra) GetBrief() extra.Data {
	return &AssetDwellingExtra{
		//BaseData:         v.BaseData,
		NumberOfBedrooms: v.NumberOfBedrooms,
		AreaSqM:          v.AreaSqM,
	}
}

func (v AssetDwellingExtra) RequiredFields() []string {
	return nil
}

func (v AssetDwellingExtra) IndexedFields() []string {
	return nil
}

func (v AssetDwellingExtra) Validate() error {
	//if err := v.BaseData.Validate(); err != nil {
	//	return err
	//}
	if v.NumberOfBedrooms < 0 {
		return validation.NewErrBadRecordFieldValue("numberOfBedrooms", "negative value")
	}
	if v.AreaSqM < 0 {
		return validation.NewErrBadRecordFieldValue("areaSqM", "negative value")
	}
	if v.RentPrice.Value < 0 {
		return validation.NewErrBadRecordFieldValue("rent_price.value", "negative value")
	}
	return nil
}
