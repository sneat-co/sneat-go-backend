package models4assetus

import (
	"github.com/strongo/validation"
	"strings"
)

type WithRegNumberField struct {
	RegNumber string `json:"regNumber" firestore:"regNumber"` // intentionally not omitempty so can be used in queries
}

// Validate validates WitRegNumberField
func (v *WithRegNumberField) Validate() error {
	if strings.TrimSpace(v.RegNumber) != v.RegNumber {
		return validation.NewErrBadRecordFieldValue("regNumber", "should not have leading or trailing spaces")
	}
	return nil
}
