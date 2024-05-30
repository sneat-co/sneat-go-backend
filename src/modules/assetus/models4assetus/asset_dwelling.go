package models4assetus

import "github.com/strongo/validation"

type AssetDwellingExtra struct {
	AssetExtraBase
	Address   string `json:"address,omitempty" firestore:"address,omitempty"`
	RentPrice struct {
		Value    float64 `json:"value,omitempty" firestore:"value,omitempty"`
		Currency string  `json:"currency,omitempty" firestore:"currency,omitempty"`
	} `json:"rent_price,omitempty" firestore:"rent_price,omitempty"`
	NumberOfBedrooms int `json:"numberOfBedrooms,omitempty" firestore:"numberOfBedrooms,omitempty"`
	AreaSqM          int `json:"areaSqM,omitempty" firestore:"areaSqM,omitempty"`
}

func (v AssetDwellingExtra) Validate() error {
	if err := v.AssetExtraBase.Validate(); err != nil {
		return err
	}
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
