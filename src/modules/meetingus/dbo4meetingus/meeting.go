package dbo4meetingus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"time"
)

var _ core.Validatable = (*Meeting)(nil)

// Excluded record
type Excluded struct {
	By     dbmodels.ByUser `json:"by,omitempty" firestore:"by,omitempty"`
	Reason string          `json:"reason,omitempty" firestore:"reason,omitempty"`
}

// Validate validates record
func (v *Excluded) Validate() error {
	if err := v.By.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("by", err.Error())
	}
	return nil
}

// MeetingMemberBrief record
type MeetingMemberBrief struct {
	briefs4contactus.ContactBrief
	Excluded *Excluded `json:"excluded,omitempty" firestore:"excluded,omitempty"`
}

func (v *MeetingMemberBrief) Equal(v2 *MeetingMemberBrief) bool {
	return v == nil && v2 == nil ||
		v != nil && v2 != nil && v.ContactBrief.Equal(&v2.ContactBrief) &&
			(v.Excluded == nil && v2.Excluded == nil || v.Excluded != nil && v2.Excluded != nil && *v.Excluded == *v2.Excluded)
}

// Validate validates record
func (v *MeetingMemberBrief) Validate() error {
	if err := v.ContactBrief.Validate(); err != nil {
		return err
	}
	if v.Excluded != nil {
		if err := v.Excluded.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("excluded", err.Error())
		}
	}
	return nil
}

// MeetingInstance DTO
type MeetingInstance interface {
	BaseMeeting() *Meeting
	Validate() error
}

var _ MeetingInstance = (*Meeting)(nil)

// Meeting record
type Meeting struct {
	dbmodels.WithUserIDs
	briefs4contactus.WithMultiTeamContacts[*MeetingMemberBrief]
	Version  int        `json:"v" firestore:"v"`
	Started  *time.Time `json:"started,omitempty" firestore:"started,omitempty"`
	Finished *time.Time `json:"finished,omitempty" firestore:"finished,omitempty"`
	Timer    *Timer     `json:"timer,omitempty" firestore:"timer,omitempty"`
}

// Meeting returns *Meeting
func (v *Meeting) Meeting() *Meeting {
	return v
}

const validateErrorPrefix = "api4meetingus record: "

// Validate validates record
func (v *Meeting) Validate() error {
	if v.Timer != nil {
		if err := v.Timer.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("timer", fmt.Errorf(
				validateErrorPrefix+"invalid 'timer' field: %w", err).Error())
		}
	}
	if err := v.validateContacts(); err != nil {
		return err
	}
	if err := v.validateUserIDs(); err != nil {
		return err
	}
	return nil
}

func (v *Meeting) validateContacts() error { // TODO: Should this be moved into contactus module?
	if usersCount := len(v.UserIDs); len(v.Contacts) == 0 && usersCount > 0 {
		// Originally thought not to be a problem. As can have no participant but granted access to spectators.
		// But decided such users should have be listed in contacts with spectator role.
		return validation.NewErrBadRecordFieldValue("contacts",
			fmt.Sprintf("meeting has no contacts but has %d entries in `userIDs` field", usersCount))
	}
	for i, member := range v.Contacts {
		newMemberErr := func(message string) error {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("members[%v]", i), message)
		}
		if member.UserID != "" {
			if !v.HasUserID(member.UserID) {
				return newMemberErr(fmt.Sprintf("reference to unknown user id: [%v]", member.UserID))
			}
		}
		if err := member.Validate(); err != nil {
			return newMemberErr(err.Error())
		}
	}
	return nil
}

func (v *Meeting) validateUserIDs() error {
	type memberRoles struct {
		isParticipant bool
		isSpectator   bool
		isExcluded    bool
		roles         []string
	}
	getRoles := func(uid string) memberRoles {
		for _, m := range v.Contacts {
			if m.UserID == uid {
				return memberRoles{
					isParticipant: m.HasRole(const4contactus.TeamMemberRoleContributor),
					isSpectator:   m.HasRole(const4contactus.TeamMemberRoleSpectator),
					isExcluded:    m.HasRole(const4contactus.TeamMemberRoleExcluded),
					roles:         m.Roles,
				}
			}
		}
		return memberRoles{}
	}
	for _, uid := range v.UserIDs {
		user := getRoles(uid)
		if !user.isParticipant && !user.isSpectator && !user.isExcluded {
			return validation.NewErrBadRecordFieldValue("userIDs",
				fmt.Sprintf(validateErrorPrefix+"user {id=%v} is neither participant, spectator or excluded, assigned roles: %v", uid, user.roles))
		}
		if user.isParticipant && user.isSpectator {
			return validation.NewErrBadRecordFieldValue("userIDs",
				fmt.Sprintf(validateErrorPrefix+"user {id=%v} is both participant and spectator", uid))
		}
	}
	return nil
}

// BaseMeeting returns base api4meetingus data
func (v *Meeting) BaseMeeting() *Meeting {
	return v
}

// HasUserID validates if api4meetingus has a user
func (v *Meeting) HasUserID(uid string) bool {
	return slice.Index(v.UserIDs, uid) >= 0
}
