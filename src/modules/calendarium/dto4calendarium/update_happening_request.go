package dto4calendarium

import "github.com/strongo/validation"

type UpdateHappeningRequest struct {
	HappeningRequest
	Title       string `json:"title"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
}

func (v UpdateHappeningRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	if len(v.Title) > 100 {
		return validation.NewErrBadRequestFieldValue("title", "too long")
	}
	if len(v.Summary) > 200 {
		return validation.NewErrBadRequestFieldValue("summary", "too long")
	}
	if len(v.Description) > 5000 {
		return validation.NewErrBadRequestFieldValue("description", "too long")
	}
	return nil
}
