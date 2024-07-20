package dto4spaceus

import (
	"fmt"
	"github.com/strongo/validation"
	"strings"
)

// SpaceMetricsRequest request
type SpaceMetricsRequest struct {
	SpaceRequest
	Metrics []string `json:"metrics"`
}

// Validate validates request
func (v *SpaceMetricsRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if len(v.Metrics) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("metrics")
	}
	for i, m := range v.Metrics {
		if strings.TrimSpace(m) == "" {
			return validation.NewErrBadRequestFieldValue("metrics", fmt.Sprintf("empty value: metrics[%d]", i))
		}
	}
	return nil
}
