package models4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
)

// ScrumSettings record
type ScrumSettings struct {
	Duration *models4teamus.MeetingDurationSettings `json:"duration" firestore:"duration"`
}

// Validate validates record
func (v *ScrumSettings) Validate() error {
	if v.Duration != nil {
		if err := v.Duration.Validate(); err != nil {
			return nil
		}
	}
	return nil
}
