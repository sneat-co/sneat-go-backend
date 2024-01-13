package models4retrospectus

import "github.com/strongo/validation"

// RetroSettings record
type RetroSettings struct {
	MaxVotesPerUser int `json:"maxVotesPerUser" firestore:"maxVotesPerUser"`
}

// Validate validates record
func (v *RetroSettings) Validate() error {
	if v == nil {
		return nil
	}
	if v.MaxVotesPerUser == 0 {
		return validation.NewErrRecordIsMissingRequiredField("maxVotesPerUser")
	}
	return nil
}
