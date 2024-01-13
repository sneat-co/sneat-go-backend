package models4assetus

import (
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

var _ AssetMain = (*DocumentMainData)(nil)

type DocumentMainData struct {
	*AssetMainDto
	*DocData // Should we move it to AssetBaseDto?
}

func (v *DocumentMainData) Validate() error {
	if err := v.AssetMainDto.Validate(); err != nil {
		return err
	}
	if err := v.DocData.Validate(); err != nil {
		return err
	}
	return nil

}
func (v *DocumentMainData) SpecificData() AssetSpecificData {
	return v.DocData
}

func (v *DocumentMainData) SetSpecificData(data AssetSpecificData) {
	v.DocData = data.(*DocData)
}

var _ AssetDbData = (*DocumentDbData)(nil)

func NewDocumentDbData() *DocumentDbData {
	return &DocumentDbData{
		DocumentMainData: new(DocumentMainData),
		AssetExtraDto:    new(AssetExtraDto),
	}
}

// DocumentDbData DTO
type DocumentDbData struct {
	*DocumentMainData
	*AssetExtraDto
}

// Validate returns error if not valid
func (v DocumentDbData) Validate() error {
	if err := v.DocumentMainData.Validate(); err != nil {
		return err
	}
	if err := v.AssetExtraDto.Validate(); err != nil {
		return err
	}
	return nil
}

type DocData struct {
	IssuedOn      string `json:"issuedOn,omitempty" firestore:"issuedOn,omitempty"`
	EffectiveFrom string `json:"effectiveFrom,omitempty" firestore:"effectiveFrom,omitempty"`
	ExpiresOn     string `json:"expiresOn,omitempty" firestore:"expiresOn,omitempty"`
}

func (v *DocData) Validate() error {
	if v.IssuedOn == "" {
		if _, err := validate.DateString(v.IssuedOn); err != nil {
			return validation.NewErrBadRecordFieldValue("issuedOn", err.Error())
		}
	}
	if v.EffectiveFrom == "" {
		if _, err := validate.DateString(v.EffectiveFrom); err != nil {
			return validation.NewErrBadRecordFieldValue("issuedOn", err.Error())
		}
	}
	if v.ExpiresOn != "" {
		if _, err := validate.DateString(v.ExpiresOn); err != nil {
			return validation.NewErrBadRecordFieldValue("issuedOn", err.Error())
		}
	}
	return nil
}
