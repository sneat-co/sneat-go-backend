package dbo4spaceus

import (
	"errors"
	"strings"
)

func ValidateShippingPointID(v string) error {
	if trimmed := strings.TrimSpace(v); trimmed != v {
		return errors.New("should not contain leading or trailing spaces")
	}
	return nil
}
