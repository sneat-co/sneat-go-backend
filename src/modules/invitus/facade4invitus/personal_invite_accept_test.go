package facade4invitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/golang/mock/gomock"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"testing"
	"time"
)

func TestAcceptPersonalInvite(t *testing.T) {
	type args struct {
		ctx         context.Context
		userContext facade.User
		request     AcceptPersonalInviteRequest
	}
	ctx := context.Background()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "nil_params",
			args:    args{ctx: ctx},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AcceptPersonalInvite(tt.args.ctx, tt.args.userContext, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("AcceptPersonalInvite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAcceptPersonalInviteRequest_Validate(t *testing.T) {
	type fields struct {
		InviteRequest InviteRequest
		Member        dbmodels.DtoWithID[*briefs4contactus.ContactBase]
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "should_return_error_for_empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &AcceptPersonalInviteRequest{
				InviteRequest: tt.fields.InviteRequest,
				Member:        tt.fields.Member,
			}
			if err := v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createOrUpdateUserRecord(t *testing.T) {
	ctx := context.Background()
	type args struct {
		user              dbo4userus.UserEntry
		userRecordError   error
		teamRecordError   error
		inviteRecordError error
		request           AcceptPersonalInviteRequest
		team              dal4teamus.TeamEntry
		teamMember        dbmodels.DtoWithID[*briefs4contactus.ContactBase]
		invite            PersonalInviteEntry
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil_params",
			args: args{
				user:            dbo4userus.NewUserEntry("test_user_id"),
				userRecordError: dal.ErrRecordNotFound,
				team: dal4teamus.NewTeamEntryWithDto("testteamid", &dbo4teamus.TeamDbo{
					TeamBrief: dbo4teamus.TeamBrief{
						RequiredCountryID: with.RequiredCountryID{
							CountryID: with.UnknownCountryID,
						},
						Type:  "family",
						Title: "Family",
					},
				}),
				teamMember: dbmodels.DtoWithID[*briefs4contactus.ContactBase]{
					ID: "test_member_id2",
					Data: &briefs4contactus.ContactBase{
						ContactBrief: briefs4contactus.ContactBrief{
							Type:   briefs4contactus.ContactTypePerson,
							Gender: "unknown",
							Names: &person.NameFields{
								FirstName: "First",
							},
							//Status:   "active",
							AgeGroup: "unknown",
						},
						//WithRequiredCountryID: dbmodels.WithRequiredCountryID{
					},
				},
				invite: NewPersonalInviteEntryWithDto("test_personal_invite_id", &dbo4invitus.PersonalInviteDbo{
					InviteDbo: dbo4invitus.InviteDbo{
						Roles: []string{"contributor"},
					},
				}),
				request: AcceptPersonalInviteRequest{
					RemoteClient: dbmodels.RemoteClientInfo{
						HostOrApp:  "unit-test",
						RemoteAddr: "localhost",
					},
					InviteRequest: InviteRequest{
						TeamRequest: dto4teamus.TeamRequest{
							TeamID: "testteamid",
						},
						InviteID: "test_personal_invite_id",
						Pin:      "1234",
					},
					Member: dbmodels.DtoWithID[*briefs4contactus.ContactBase]{
						ID: "test_member_id",
						Data: &briefs4contactus.ContactBase{
							ContactBrief: briefs4contactus.ContactBrief{
								Type:     briefs4contactus.ContactTypePerson,
								Gender:   "unknown",
								AgeGroup: "unknown",
							},
							Status: "active",
							//WithRequiredCountryID: dbmodels.WithRequiredCountryID{
							//	CountryID: dbmodels.UnknownCountryID,
							//},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			//
			tt.args.user.Record.SetError(tt.args.userRecordError)
			tt.args.team.Record.SetError(tt.args.teamRecordError)
			tt.args.invite.Record.SetError(tt.args.inviteRecordError)
			//
			tx := mocks4dal.NewMockReadwriteTransaction(mockCtrl)
			if tt.args.userRecordError == nil && tt.args.teamRecordError == nil && tt.args.inviteRecordError == nil {
				tx.EXPECT().Insert(gomock.Any(), tt.args.user.Record).Return(nil)
			}
			now := time.Now()
			params := dal4contactus.NewContactusTeamWorkerParams(tt.args.user.ID, tt.args.team.ID)
			if err := createOrUpdateUserRecord(ctx, tx, now, tt.args.user, tt.args.request, params, tt.args.teamMember.Data, tt.args.invite); err != nil {
				if !tt.wantErr {
					t.Errorf("createOrUpdateUserRecord() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			userDto := tt.args.user.Data
			assert.Equal(t, now, userDto.CreatedAt, "CreatedAt")
			assert.Equal(t, tt.args.request.Member.Data.Gender, userDto.Gender, "Gender")
			assert.Equal(t, 1, len(userDto.Teams), "len(Teams)")
			assert.Equal(t, 1, len(userDto.TeamIDs), "len(TeamIDs)")
			assert.True(t, slice.Contains(userDto.TeamIDs, tt.args.request.TeamID), "TeamIDs contains tt.args.request.TeamID")
			teamBrief := userDto.Teams[tt.args.request.TeamID]
			assert.NotNil(t, teamBrief, "Teams[tt.args.request.TeamID]")
		})
	}
}

func Test_updateInviteRecord(t *testing.T) {
	ctx := context.Background()
	type args struct {
		uid    string
		invite PersonalInviteEntry
		status string
	}
	now := time.Now()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should_pass",
			args: args{
				status: "accepted",
				invite: NewPersonalInviteEntryWithDto("test_invite_id1", &dbo4invitus.PersonalInviteDbo{
					ToTeamMemberID: "to_member_id2",
					Address:        "to.test.user@example.com",
					InviteDbo: dbo4invitus.InviteDbo{
						Pin:    "1234",
						TeamID: "testteamid1",
						Team: dbo4invitus.InviteTeam{
							ID:    "testteamid1",
							Type:  "family",
							Title: "Family",
						},
						CreatedAt: time.Now(),
						Created: dbmodels.CreatedInfo{
							Client: dbmodels.RemoteClientInfo{
								HostOrApp:  "unit-test",
								RemoteAddr: "127.0.0.1",
							},
						},
						InviteBase: dbo4invitus.InviteBase{
							Type:    "personal",
							Channel: "email",
							From: dbo4invitus.InviteFrom{
								InviteContact: dbo4invitus.InviteContact{
									UserID:   "from_user_id1",
									MemberID: "from_member_id1",
									Title:    "From ID 1",
								},
							},
							To: &dbo4invitus.InviteTo{
								InviteContact: dbo4invitus.InviteContact{
									Title:    "To ID 2",
									MemberID: "to_member_id2",
									Channel:  "email",
									Address:  "to.test.user@example.com",
								},
							},
						},
						Roles: []string{"contributor"},
					},
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			tx := mocks4dal.NewMockReadwriteTransaction(mockCtrl)
			tx.EXPECT().Update(ctx, tt.args.invite.Key, gomock.Any()).Return(nil)
			assert.Equal(t, "", tt.args.invite.Data.To.UserID)
			if err := updateInviteRecord(ctx, tx, tt.args.uid, now, tt.args.invite, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("updateInviteRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.args.status, tt.args.invite.Data.Status)
			assert.Equal(t, tt.args.uid, tt.args.invite.Data.To.UserID)
		})
	}
}

func Test_updateTeamRecord(t *testing.T) {
	type args struct {
		uid           string
		memberID      string
		team          dal4teamus.TeamEntry
		contactusTeam dal4contactus.ContactusTeamModuleEntry
		requestMember dbmodels.DtoWithID[*briefs4contactus.ContactBase]
	}
	testMember := dbmodels.DtoWithID[*briefs4contactus.ContactBase]{
		ID:   "test_member_id1",
		Data: &briefs4contactus.ContactBase{},
	}
	tests := []struct {
		name           string
		teamRecordErr  error
		args           args
		wantTeamMember dbmodels.DtoWithID[*briefs4contactus.ContactBase]
		wantErr        bool
	}{
		{
			name:          "should_pass",
			teamRecordErr: nil,
			args: args{
				uid:      "test_user_id",
				memberID: "test_member_id1",
				team: dal4teamus.NewTeamEntryWithDto("testteamid", &dbo4teamus.TeamDbo{
					TeamBrief: dbo4teamus.TeamBrief{
						Type:  "family",
						Title: "Family",
					},
				}),
				contactusTeam: dal4contactus.NewContactusTeamModuleEntryWithData("testteamid", &models4contactus.ContactusTeamDbo{
					WithSingleTeamContactsWithoutContactIDs: briefs4contactus.WithSingleTeamContactsWithoutContactIDs[*briefs4contactus.ContactBrief]{
						WithContactsBase: briefs4contactus.WithContactsBase[*briefs4contactus.ContactBrief]{
							WithContactBriefs: briefs4contactus.WithContactBriefs[*briefs4contactus.ContactBrief]{
								Contacts: map[string]*briefs4contactus.ContactBrief{
									testMember.ID: &testMember.Data.ContactBrief,
								},
							},
						},
					},
				}),
				requestMember: dbmodels.DtoWithID[*briefs4contactus.ContactBase]{
					ID: testMember.ID,
					Data: &briefs4contactus.ContactBase{
						ContactBrief: briefs4contactus.ContactBrief{
							Names: &person.NameFields{
								FirstName: "First name",
							},
						},
					},
				},
			},
			wantErr:        false,
			wantTeamMember: testMember,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//mockCtrl := gomock.NewController(t)
			//tx := mocks4dal.NewMockReadwriteTransaction(mockCtrl)
			//tx.EXPECT().Update(gomock.Any(), tt.args.team.Key, gomock.Any()).Return(nil)
			//tx.EXPECT().Update(gomock.Any(), tt.args.contactusTeam.Key, gomock.Any()).Return(nil)
			tt.args.contactusTeam.Record.SetError(tt.teamRecordErr)
			params := dal4contactus.NewContactusTeamWorkerParams(tt.args.uid, tt.args.team.ID)
			params.TeamModuleEntry.Data.AddContact(tt.args.memberID, &tt.args.requestMember.Data.ContactBrief)
			params.TeamModuleEntry.Data.AddUserID(tt.args.uid)
			params.Team.Data.AddUserID(tt.args.uid)
			gotTeamMember, err := updateTeamRecord(tt.args.uid, tt.args.memberID, params, tt.args.requestMember)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateTeamRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, gotTeamMember, "gotTeamMember is nil")
			//if !reflect.DeepEqual(gotTeamMember, tt.wantTeamMember) {
			//	t.Errorf("updateTeamRecord() gotTeamMember = %v, want %v", gotTeamMember, tt.wantTeamMember)
			//}
		})
	}
}
