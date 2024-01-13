package dto4assetus

import (
	"github.com/strongo/validation"
	"strings"
)

type UpdateAssetRequest struct {
	AssetRequest
	RegNumber *string `json:"regNumber,omitempty"`
}

func (v UpdateAssetRequest) Validate() error {
	if err := v.AssetRequest.Validate(); err != nil {
		return err
	}
	if v.RegNumber != nil {
		regNumber := *v.RegNumber
		if strings.TrimSpace(regNumber) != regNumber {
			return validation.NewErrBadRequestFieldValue("regNumber", "should not have leading or trailing spaces")
		}
	}
	return nil
}
