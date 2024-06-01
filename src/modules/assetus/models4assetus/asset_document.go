package models4assetus

import (
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"time"
)

func init() {
	RegisterAssetExtraFactory(AssetExtraTypeDocument, func() AssetExtra {
		return new(AssetDocumentExtra)
	})
}

var _ AssetExtra = (*AssetDocumentExtra)(nil)

type AssetDocumentExtra struct {
	AssetExtraBase
	WithRegNumberField
	IssuedOn      string `json:"issuedOn,omitempty" firestore:"issuedOn,omitempty"`
	EffectiveFrom string `json:"effectiveFrom,omitempty" firestore:"effectiveFrom,omitempty"`
	ExpiresOn     string `json:"expiresOn,omitempty" firestore:"expiresOn,omitempty"`
}

func (v *AssetDocumentExtra) RequiredFields() []string {
	return []string{""}
}

func (v *AssetDocumentExtra) IndexedFields() []string {
	return []string{"expiresOn", "effectiveFrom"}
}

func (v *AssetDocumentExtra) GetBrief() AssetExtra {
	return &AssetDocumentExtra{
		AssetExtraBase:     v.AssetExtraBase,
		IssuedOn:           v.IssuedOn,
		EffectiveFrom:      v.EffectiveFrom,
		ExpiresOn:          v.ExpiresOn,
		WithRegNumberField: v.WithRegNumberField,
	}
}

func (v *AssetDocumentExtra) Validate() (err error) {
	if err := v.AssetExtraBase.Validate(); err != nil {
		return err
	}
	if err := v.WithRegNumberField.Validate(); err != nil {
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
