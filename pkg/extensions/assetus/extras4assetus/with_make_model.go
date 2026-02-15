package extras4assetus

import (
	"strings"

	"github.com/strongo/validation"
)

type WithMakeModelFields struct {
	Make  string `json:"make,omitempty" firestore:"make,omitempty"`
	Model string `json:"model,omitempty" firestore:"model,omitempty"`
}

// GenerateTitleFromMakeModelAndRegNumber generates asset title from vehicle data
func (v *WithMakeModelFields) GenerateTitleFromMakeModelAndRegNumber(reNumber string) string {
	title := make([]string, 0, 4)
	if v.Make != "" {
		title = append(title, v.Make)
	}
	if v.Model != "" {
		title = append(title, v.Model)
	}
	if reNumber != "" {
		title = append(title, "#", reNumber)
	}
	if len(title) == 0 {
		return ""
	}
	return strings.Join(title, " ")
}

func (v *WithMakeModelFields) Validate() error {
	if makeValue := strings.TrimSpace(v.Make); makeValue == "" {
		return validation.NewErrRecordIsMissingRequiredField("make")
	} else if makeValue != v.Make {
		return validation.NewErrBadRecordFieldValue("make", "should not have leading or trailing spaces")
	}
	if model := strings.TrimSpace(v.Model); model == "" {
		return validation.NewErrRecordIsMissingRequiredField("model")
	} else if model != v.Model {
		return validation.NewErrBadRecordFieldValue("model", "should not have leading or trailing spaces")
	}
	return nil
}

type WithMakeModelRegNumberFields struct {
	WithMakeModelFields
	WithOptionalRegNumberField
}

func (v *WithMakeModelRegNumberFields) Validate() error {
	if err := v.WithMakeModelFields.Validate(); err != nil {
		return err
	}
	if err := v.WithOptionalRegNumberField.Validate(); err != nil {
		return err
	}
	return nil
}
