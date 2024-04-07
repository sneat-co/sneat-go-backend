package dto4logist

import (
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strings"
)

type SetFieldsRequest struct {
	SetDates   map[string]string `json:"setDates,omitempty"`
	SetStrings map[string]string `json:"setStrings,omitempty"`
}

func (v SetFieldsRequest) Validate() error {
	for name, value := range v.SetDates {
		if _, err := validate.DateString(value); err != nil {
			return validation.NewErrBadRequestFieldValue("setDates."+name, err.Error())
		}
	}
	for name, value := range v.SetStrings {
		if strings.TrimSpace(value) != value {
			return validation.NewErrBadRequestFieldValue("setStrings."+name, "must not start or end with space")
		}
	}
	return nil
}
