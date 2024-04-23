package models4assetus

import (
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"time"
)

type AssetDocumentExtra struct {
	IssuedOn      string `json:"issuedOn,omitempty" firestore:"issuedOn,omitempty"`
	EffectiveFrom string `json:"effectiveFrom,omitempty" firestore:"effectiveFrom,omitempty"`
	ExpiresOn     string `json:"expiresOn,omitempty" firestore:"expiresOn,omitempty"`
}

func (v *AssetDocumentExtra) Validate() (err error) {
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
