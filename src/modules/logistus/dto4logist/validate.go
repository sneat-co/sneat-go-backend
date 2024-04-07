package dto4logist

import (
	"github.com/strongo/validation"
	"strings"
)

func validateID(field, v string) error {
	if strings.TrimSpace(v) == "" {
		return validation.NewErrRequestIsMissingRequiredField(field)
	}
	if strings.ContainsAny(v, " \t\r") {
		return validation.NewErrBadRequestFieldValue(field, "must not contain spaces")
	}
	return nil
}

func validateContainerID(field, v string) error {
	return validateID(field, v)
}
