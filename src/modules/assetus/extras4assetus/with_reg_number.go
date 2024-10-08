package extras4assetus

import (
	"github.com/strongo/validation"
	"strings"
)

type WithOptionalRegNumberField struct {
	RegNumber string `json:"regNumber,omitempty" firestore:"regNumber,omitempty" dalgo:"regNumber,omitempty"`
}

// Validate validates WitRegNumberField
func (v *WithOptionalRegNumberField) Validate() error {
	if regNumber := strings.TrimSpace(v.RegNumber); regNumber != v.RegNumber {
		return validation.NewErrBadRecordFieldValue("regNumber", "should not have leading or trailing spaces")
	}
	return nil
}
