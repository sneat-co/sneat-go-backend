package models4assetus

import "github.com/strongo/validation"

type DwellingData struct {
	BedRooms  int `json:"bedRooms,omitempty" firestore:"bedRooms,omitempty"`
	BathRooms int `json:"bathRooms,omitempty" firestore:"bathRooms,omitempty"`
}

func (v DwellingData) Validate() error {
	if v.BedRooms < 0 {
		return validation.NewErrBadRecordFieldValue("bedRooms", "negative value")
	}
	if v.BathRooms < 0 {
		return validation.NewErrBadRecordFieldValue("bathRooms", "negative value")
	}
	return nil
}

var _ AssetMain = (*AssetDtoDwelling)(nil)

type AssetDtoDwelling struct {
	AssetMainDto
	DwellingData
}

func (v *AssetDtoDwelling) SpecificData() AssetSpecificData {
	return &v.DwellingData
}

func (v *AssetDtoDwelling) SetSpecificData(data AssetSpecificData) {
	v.DwellingData = data.(DwellingData)
}

func (v *AssetDtoDwelling) Validate() error {
	if err := v.AssetMainDto.Validate(); err != nil {
		return err
	}
	if err := v.DwellingData.Validate(); err != nil {
		return err
	}
	return nil
}

var _ AssetDbData = (*DwellingAssetDbData)(nil)

func NewDwellingAssetDbData() *DwellingAssetDbData {
	return &DwellingAssetDbData{
		AssetDtoDwelling: new(AssetDtoDwelling),
		AssetExtraDto:    new(AssetExtraDto),
	}
}

type DwellingAssetDbData struct {
	*AssetDtoDwelling
	*AssetExtraDto
}

func (v DwellingAssetDbData) Validate() error {
	if err := v.AssetDtoDwelling.Validate(); err != nil {
		return err
	}
	if err := v.AssetExtraDto.Validate(); err != nil {
		return err
	}
	return nil
}
