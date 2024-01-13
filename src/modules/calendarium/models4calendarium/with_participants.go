package models4calendarium

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

type WithParticipants struct {
	// Participants keeps contact info specific to the happening.
	// Map key is expected to be valid dbmodels.TeamItemID to support contacts from multiple teams.
	Participants map[string]*HappeningParticipant `json:"participants,omitempty" firestore:"participants,omitempty"`
}

func (v *WithParticipants) AddParticipant(teamID, contactID string, participant *HappeningParticipant) []dal.Update {
	id := dbmodels.NewTeamItemID(teamID, contactID)
	if participant == nil {
		participant = &HappeningParticipant{}
	}
	if v.Participants == nil {
		v.Participants = make(map[string]*HappeningParticipant)
	}
	if _, ok := v.Participants[string(id)]; ok {
		return []dal.Update{}
	}
	v.Participants[string(id)] = participant
	return []dal.Update{
		{
			Field: "participants." + string(id),
			Value: participant,
		},
	}
}

func (v *WithParticipants) RemoveParticipant(teamID, contactID string) []dal.Update {
	id := dbmodels.NewTeamItemID(teamID, contactID)
	if v.Participants == nil {
		return []dal.Update{}
	}
	if _, ok := v.Participants[string(id)]; !ok {
		return []dal.Update{}
	}
	delete(v.Participants, string(id))
	return []dal.Update{
		{
			Field: "participants." + string(id),
			Value: dal.DeleteField,
		},
	}
}

func (v *WithParticipants) Validate() error {
	for contactID, participant := range v.Participants {
		if contactID == "" {
			return validation.NewErrBadRecordFieldValue("participants", "contactID is empty")
		}
		field := func() string {
			return fmt.Sprintf("participants[%s]", contactID)
		}
		if err := dbmodels.TeamItemID(contactID).Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(field(), err.Error())
		}
		if participant != nil {
			if err := participant.Validate(); err != nil {
				return validation.NewErrBadRecordFieldValue(field(), err.Error())
			}
		}
	}
	return nil
}

func (v *WithParticipants) Updates() (updates []dal.Update) {
	if len(v.Participants) == 0 {
		updates = append(updates, dal.Update{
			Field: "participants",
			Value: dal.DeleteField,
		})
		return
	}
	updates = append(updates, dal.Update{
		Field: "participants",
		Value: v.Participants,
	})
	return
}
