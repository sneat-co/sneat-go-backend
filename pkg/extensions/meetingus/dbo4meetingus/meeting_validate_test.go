package dbo4meetingus

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

// newMemberBrief returns a minimally valid MeetingMemberBrief (person with a title).
func newMemberBrief(userID string, roles ...string) *MeetingMemberBrief {
	m := &MeetingMemberBrief{}
	m.Type = briefs4contactus.ContactTypePerson
	m.Title = "Member " + userID
	m.UserID = userID
	m.Roles = roles
	return m
}

func TestExcluded_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       Excluded
		wantErr bool
	}{
		{"valid", Excluded{By: dbmodels.ByUser{UID: "u1"}, Reason: "spam"}, false},
		{"missing_by", Excluded{Reason: "spam"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMeetingMemberBrief_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       *MeetingMemberBrief
		wantErr bool
	}{
		{"valid", newMemberBrief("u1"), false},
		{"invalid_contact_brief", &MeetingMemberBrief{}, true}, // missing type
		{
			name: "valid_with_excluded",
			v: func() *MeetingMemberBrief {
				m := newMemberBrief("u1")
				m.Excluded = &Excluded{By: dbmodels.ByUser{UID: "u2"}}
				return m
			}(),
			wantErr: false,
		},
		{
			name: "invalid_excluded",
			v: func() *MeetingMemberBrief {
				m := newMemberBrief("u1")
				m.Excluded = &Excluded{} // missing By
				return m
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMeetingMemberBrief_Equal(t *testing.T) {
	base := newMemberBrief("u1")
	same := newMemberBrief("u1")
	other := newMemberBrief("u2")

	excludedA := newMemberBrief("u1")
	excludedA.Excluded = &Excluded{By: dbmodels.ByUser{UID: "u9"}}
	excludedB := newMemberBrief("u1")
	excludedB.Excluded = &Excluded{By: dbmodels.ByUser{UID: "u9"}}

	tests := []struct {
		name string
		a    *MeetingMemberBrief
		b    *MeetingMemberBrief
		want bool
	}{
		{"both_nil", nil, nil, true},
		{"one_nil", base, nil, false},
		{"equal", base, same, true},
		{"different_user", base, other, false},
		{"equal_with_excluded", excludedA, excludedB, true},
		{"one_excluded_one_not", base, excludedA, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.a.Equal(tt.b))
		})
	}
}

func TestMeeting_Meeting_and_BaseMeeting(t *testing.T) {
	m := &Meeting{Version: 2}
	assert.Same(t, m, m.Meeting())
	assert.Same(t, m, m.BaseMeeting())
}

func TestMeeting_Validate_Contacts_and_UserIDs(t *testing.T) {
	// helper to build a meeting with given userIDs and contacts map
	build := func(userIDs []string, contacts map[string]*MeetingMemberBrief) *Meeting {
		m := &Meeting{}
		m.UserIDs = userIDs
		m.Contacts = contacts
		return m
	}

	tests := []struct {
		name    string
		meeting *Meeting
		wantErr bool
	}{
		{
			name:    "empty",
			meeting: &Meeting{},
			wantErr: false,
		},
		{
			name:    "users_without_contacts",
			meeting: build([]string{"u1"}, nil),
			wantErr: true,
		},
		{
			name: "valid_contributor",
			meeting: build([]string{"u1"}, map[string]*MeetingMemberBrief{
				"c1": newMemberBrief("u1", const4contactus.SpaceMemberRoleContributor),
			}),
			wantErr: false,
		},
		{
			name: "valid_spectator",
			meeting: build([]string{"u1"}, map[string]*MeetingMemberBrief{
				"c1": newMemberBrief("u1", const4contactus.SpaceMemberRoleSpectator),
			}),
			wantErr: false,
		},
		{
			name: "contact_references_unknown_user",
			meeting: build([]string{"u1"}, map[string]*MeetingMemberBrief{
				"c1": newMemberBrief("uX", const4contactus.SpaceMemberRoleContributor),
			}),
			wantErr: true,
		},
		{
			name: "invalid_member_brief",
			meeting: build(nil, map[string]*MeetingMemberBrief{
				"c1": {}, // missing type -> invalid
			}),
			wantErr: true,
		},
		{
			name: "user_with_no_role",
			meeting: build([]string{"u1"}, map[string]*MeetingMemberBrief{
				"c1": newMemberBrief("u1"),
			}),
			wantErr: true,
		},
		{
			name: "user_both_participant_and_spectator",
			meeting: build([]string{"u1"}, map[string]*MeetingMemberBrief{
				"c1": newMemberBrief("u1",
					const4contactus.SpaceMemberRoleContributor,
					const4contactus.SpaceMemberRoleSpectator),
			}),
			wantErr: true,
		},
		{
			name: "excluded_user_ok",
			meeting: build([]string{"u1"}, map[string]*MeetingMemberBrief{
				"c1": newMemberBrief("u1", const4contactus.SpaceMemberRoleExcluded),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.meeting.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMeeting_Validate_Timer(t *testing.T) {
	m := &Meeting{Timer: &Timer{Status: "bogus", By: dbmodels.ByUser{UID: "u1"}}}
	assert.Error(t, m.Validate())

	m2 := &Meeting{Timer: &Timer{Status: TimerStatusActive, By: dbmodels.ByUser{UID: "u1"}}}
	assert.NoError(t, m2.Validate())
}
