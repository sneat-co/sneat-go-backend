package models4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
)

// ScrumTeamDto is a DTO for scrum team
type ScrumTeamDto struct {
	ScrumSettings *ScrumSettings                 `json:"scrumSettings,omitempty" firestore:"scrumSettings,omitempty"`
	Active        *models4teamus.TeamMeetingInfo `json:"active,omitempty" firestore:"active,omitempty"`
	Last          *models4teamus.TeamMeetingInfo `json:"last,omitempty" firestore:"last,omitempty"`
}

// Validate returns error if not valid
func (v ScrumTeamDto) Validate() error {
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
