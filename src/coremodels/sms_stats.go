package coremodels

import (
	"github.com/strongo/validation"
)

type SmsStats struct {
	SmsCount int64   `firestore:"smsCount,omitempty" datastore:"smsCount,noindex,omitempty"`
	SmsCost  float64 `firestore:"smsCost,omitempty" datastore:"smsCost,noindex,omitempty"`
}

func (v *SmsStats) Validate() error {
	if v.SmsCount < 0 {
		return validation.NewErrBadRecordFieldValue("smsCount", "is negative")
	}
	if v.SmsCost < 0 {
		return validation.NewErrBadRecordFieldValue("smsCost", "is negative")
	}
	return nil
}
