package facade4invitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/golang/mock/gomock"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"slices"
	"testing"
	"time"
)

func TestAcceptPersonalInvite(t *testing.T) {
	type args struct {
		ctx     context.Context
		userCtx facade.UserContext
		request AcceptPersonalInviteRequest
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
			if err := AcceptPersonalInvite(tt.args.ctx, tt.args.userCtx, tt.args.request); (err != nil) != tt.wantErr {
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
		team              dbo4spaceus.SpaceEntry
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
				team: dbo4spaceus.NewSpaceEntryWithDbo("testteamid", &dbo4spaceus.SpaceDbo{
					SpaceBrief: dbo4spaceus.SpaceBrief{
						OptionalCountryID: with.OptionalCountryID{
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
						SpaceRequest: dto4spaceus.SpaceRequest{
							SpaceID: "testteamid",
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
			params := dal4contactus.NewContactusSpaceWorkerParams(facade.NewUserContext(tt.args.user.ID), tt.args.team.ID)
			if err := createOrUpdateUserRecord(ctx, tx, now, tt.args.user, tt.args.request, params, tt.args.teamMember.Data, tt.args.invite); err != nil {
				if !tt.wantErr {
					t.Errorf("createOrUpdateUserRecord() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			userDto := tt.args.user.Data
			assert.Equal(t, now, userDto.CreatedAt, "CreatedAt")
			assert.Equal(t, tt.args.request.Member.Data.Gender, userDto.Gender, "Gender")
			assert.Equal(t, 1, len(userDto.Spaces), "len(Spaces)")
			assert.Equal(t, 1, len(userDto.SpaceIDs), "len(SpaceIDs)")
			assert.True(t, slices.Contains(userDto.SpaceIDs, tt.args.request.SpaceID), "SpaceIDs contains tt.args.request.Space")
			teamBrief := userDto.Spaces[tt.args.request.SpaceID]
			assert.NotNil(t, teamBrief, "Spaces[tt.args.request.Space]")
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
					ToSpaceMemberID: "to_member_id2",
					Address:         "to.test.user@example.com",
					InviteDbo: dbo4invitus.InviteDbo{
						Pin:     "1234",
						SpaceID: "testteamid1",
						Space: dbo4invitus.InviteSpace{
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
									Title:    "From ContactID 1",
								},
							},
							To: &dbo4invitus.InviteTo{
								InviteContact: dbo4invitus.InviteContact{
									Title:    "To ContactID 2",
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

func Test_updateSpaceRecord(t *testing.T) {
	type args struct {
		uid            string
		memberID       string
		team           dbo4spaceus.SpaceEntry
		contactusSpace dal4contactus.ContactusSpaceEntry
		requestMember  dbmodels.DtoWithID[*briefs4contactus.ContactBase]
	}
	testMember := dbmodels.DtoWithID[*briefs4contactus.ContactBase]{
		ID:   "test_member_id1",
		Data: &briefs4contactus.ContactBase{},
	}
	tests := []struct {
		name            string
		teamRecordErr   error
		args            args
		wantSpaceMember dbmodels.DtoWithID[*briefs4contactus.ContactBase]
		wantErr         bool
	}{
		{
			name:          "should_pass",
			teamRecordErr: nil,
			args: args{
				uid:      "test_user_id",
				memberID: "test_member_id1",
				team: dbo4spaceus.NewSpaceEntryWithDbo("testteamid", &dbo4spaceus.SpaceDbo{
					SpaceBrief: dbo4spaceus.SpaceBrief{
						Type:  "family",
						Title: "Family",
					},
				}),
				contactusSpace: dal4contactus.NewContactusSpaceEntryWithData("testteamid", &dbo4contactus.ContactusSpaceDbo{
					WithSingleSpaceContactsWithoutContactIDs: briefs4contactus.WithSingleSpaceContactsWithoutContactIDs[*briefs4contactus.ContactBrief]{
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
			wantErr:         false,
			wantSpaceMember: testMember,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//mockCtrl := gomock.NewController(t)
			//tx := mocks4dal.NewMockReadwriteTransaction(mockCtrl)
			//tx.EXPECT().Update(gomock.Any(), tt.args.team.Key, gomock.Any()).Return(nil)
			//tx.EXPECT().Update(gomock.Any(), tt.args.contactusSpace.Key, gomock.Any()).Return(nil)
			tt.args.contactusSpace.Record.SetError(tt.teamRecordErr)
			params := dal4contactus.NewContactusSpaceWorkerParams(facade.NewUserContext(tt.args.uid), tt.args.team.ID)
			params.SpaceModuleEntry.Data.AddContact(tt.args.memberID, &tt.args.requestMember.Data.ContactBrief)
			params.SpaceModuleEntry.Data.AddUserID(tt.args.uid)
			params.Space.Data.AddUserID(tt.args.uid)
			gotSpaceMember, err := updateSpaceRecord(tt.args.uid, tt.args.memberID, params, tt.args.requestMember)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateSpaceRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, gotSpaceMember, "gotSpaceMember is nil")
			//if !reflect.DeepEqual(gotSpaceMember, tt.wantSpaceMember) {
			//	t.Errorf("updateSpaceRecord() gotSpaceMember = %v, want %v", gotSpaceMember, tt.wantSpaceMember)
			//}
		})
	}
}
