package extras4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/core/extra"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"time"
)

func init() {
	RegisterAssetExtraFactory(AssetExtraTypeDocument, func() briefs4assetus.AssetExtra {
		return new(AssetDocumentExtra)
	})
}

var _ extra.Data = (*AssetDocumentExtra)(nil)
var _ briefs4assetus.AssetExtra = (*AssetDocumentExtra)(nil)

type AssetDocumentExtra struct {
	WithOptionalRegNumberField
	IssuedOn      string `json:"issuedOn,omitempty" firestore:"issuedOn,omitempty"`
	EffectiveFrom string `json:"effectiveFrom,omitempty" firestore:"effectiveFrom,omitempty"`
	ExpiresOn     string `json:"expiresOn,omitempty" firestore:"expiresOn,omitempty"`
}

func (v *AssetDocumentExtra) ValidateWithAssetBrief(assetBrief briefs4assetus.AssetBrief) error {
	if err := v.Validate(); err != nil {
		return err
	}
	if assetBrief.Title == "" && v.RegNumber == "" {
		return validation.NewValidationError("document asset should have at least 1 of next fields: title, regNumber")
	}
	return nil
}

func (v *AssetDocumentExtra) RequiredFields() []string {
	return []string{""}
}

func (v *AssetDocumentExtra) IndexedFields() []string {
	return []string{"expiresOn", "effectiveFrom"}
}

func (v *AssetDocumentExtra) GetBrief() extra.Data {
	return &AssetDocumentExtra{
		//BaseData:           v.BaseData,
		IssuedOn:                   v.IssuedOn,
		EffectiveFrom:              v.EffectiveFrom,
		ExpiresOn:                  v.ExpiresOn,
		WithOptionalRegNumberField: v.WithOptionalRegNumberField,
	}
}

func (v *AssetDocumentExtra) Validate() (err error) {
	//if err := v.BaseData.Validate(); err != nil {
	//	return err
	//}
	if err := v.WithOptionalRegNumberField.Validate(); err != nil {
		return err
	}
	if v.IssuedOn != "" {
		if _, err = validate.DateString(v.IssuedOn); err != nil {
			return validation.NewErrBadRecordFieldValue("issuedOn", err.Error())
		}
	}
	var effectiveFrom, expiresOn time.Time
	if v.EffectiveFrom != "" {
		if effectiveFrom, err = validate.DateString(v.EffectiveFrom); err != nil {
			return validation.NewErrBadRecordFieldValue("effectiveFrom", err.Error())
		}
	}

	if v.ExpiresOn != "" {
		if expiresOn, err = validate.DateString(v.ExpiresOn); err != nil {
			return validation.NewErrBadRecordFieldValue("issuedOn", err.Error())
		}
	}
	if !effectiveFrom.IsZero() && !expiresOn.IsZero() || expiresOn.Before(effectiveFrom) {
		return validation.NewErrBadRecordFieldValue("expiresOn", "is before `effectiveFrom`")
	}
	return nil
}
