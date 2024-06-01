package models4assetus

import (
	"github.com/strongo/validation"
	"strings"
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
	if strings.TrimSpace(v.Make) == "" {
		return validation.NewErrRecordIsMissingRequiredField("make")
	}
	if strings.TrimSpace(v.Model) == "" {
		return validation.NewErrRecordIsMissingRequiredField("model")
	}
	return nil
}

type WithMakeModelRegNumberFields struct {
	WithMakeModelFields
	WithRegNumberField
}

func (v *WithMakeModelRegNumberFields) Validate() error {
	if err := v.WithMakeModelFields.Validate(); err != nil {
		return err
	}
	if err := v.WithRegNumberField.Validate(); err != nil {
		return err
	}
	return nil
}
