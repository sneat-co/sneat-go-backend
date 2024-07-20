package dbo4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
)

// ScrumSpaceDto is a DTO for scrum team
type ScrumSpaceDto struct {
	ScrumSettings *ScrumSettings                `json:"scrumSettings,omitempty" firestore:"scrumSettings,omitempty"`
	Active        *dbo4spaceus.SpaceMeetingInfo `json:"active,omitempty" firestore:"active,omitempty"`
	Last          *dbo4spaceus.SpaceMeetingInfo `json:"last,omitempty" firestore:"last,omitempty"`
}

// Validate returns error if not valid
func (v ScrumSpaceDto) Validate() error {
	if v.ScrumSettings != nil {
		if err := v.ScrumSettings.Validate(); err != nil {
			return err
		}
	}
	if v.Active != nil {
		if err := v.Active.Validate(); err != nil {
			return err
		}
	}
	if v.Last != nil {
		if err := v.Last.Validate(); err != nil {
			return err
		}
	}
	return nil
}
