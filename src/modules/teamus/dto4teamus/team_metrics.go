package dto4teamus

import (
	"fmt"
	"github.com/strongo/validation"
	"strings"
)

// TeamMetricsRequest request
type TeamMetricsRequest struct {
	TeamRequest
	Metrics []string `json:"metrics"`
}

// Validate validates request
func (v *TeamMetricsRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if len(v.Metrics) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("metrics")
	}
	for i, m := range v.Metrics {
		if strings.TrimSpace(m) == "" {
			return validation.NewErrBadRequestFieldValue("metrics", fmt.Sprintf("empty value: metrics[%v]", i))
		}
	}
	return nil
}
