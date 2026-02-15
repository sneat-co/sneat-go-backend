package dbo4scrumus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
)

// ScrumSettings record
type ScrumSettings struct {
	Duration *dbo4spaceus.MeetingDurationSettings `json:"duration" firestore:"duration"`
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
