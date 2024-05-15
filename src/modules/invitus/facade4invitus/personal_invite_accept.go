package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// AcceptPersonalInviteRequest holds parameters for accepting a personal invite
type AcceptPersonalInviteRequest struct {
	InviteRequest
	RemoteClient dbmodels.RemoteClientInfo `json:"remoteClient"`
	MemberID     string
	Member       dbmodels.DtoWithID[*briefs4contactus.ContactBase] `json:"member"`
	//FullName string                      `json:"fullName"`
	//Email    string                      `json:"email"`
}

// Validate validates request
func (v *AcceptPersonalInviteRequest) Validate() error {
	if err := v.InviteRequest.Validate(); err != nil {
		return err
	}
	if err := v.Member.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("member", err.Error())
	}
	//if v.FullName == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("FullName")
	//}
	//if v.Email == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("Email")
	//}
	return nil
}

// AcceptPersonalInvite accepts personal invite and joins user to a team.
// If needed a user record should be created
func AcceptPersonalInvite(ctx context.Context, userContext facade.User, request AcceptPersonalInviteRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	uid := userContext.GetID()

	return dal4contactus.RunContactusTeamWorker(ctx, userContext, request.TeamRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusTeamWorkerParams) error {
			invite, member, err := getPersonalInviteRecords(ctx, tx, params, request.InviteID, request.Member.ID)
			if err != nil {
				return err
			}
			if invite.Dto.Status != "active" {
				return fmt.Errorf("invite status is not equal to 'active', got '%s'", invite.Dto.Status)
			}

			if invite.Dto.Pin != request.Pin {
				return fmt.Errorf("%w: pin code does not match", facade.ErrBadRequest)
			}

			user := models4userus.NewUserContext(uid)
			if err = facade4userus.GetUserByID(ctx, tx, user.Record); err != nil {
				if !dal.IsNotFound(err) {
					return err
				}
			}

			now := time.Now()

			if err = updateInviteRecord(ctx, tx, uid, now, invite, "accepted"); err != nil {
				return fmt.Errorf("failed to update invite record: %w", err)
			}

			var teamMember *briefs4contactus.ContactBase
			if teamMember, err = updateTeamRecord(uid, invite.Dto.ToTeamMemberID, params, request.Member); err != nil {
				return fmt.Errorf("failed to update team record: %w", err)
			}

			memberContext := dal4contactus.NewContactEntry(params.Team.ID, member.ID)

			if err = updateMemberRecord(ctx, tx, uid, memberContext, request.Member.Data, teamMember); err != nil {
				return fmt.Errorf("failed to update team member record: %w", err)
			}

			if err = createOrUpdateUserRecord(ctx, tx, now, user, request, params, teamMember, invite); err != nil {
				return fmt.Errorf("failed to create or update user record: %w", err)
			}

			return err
		})
}

func updateInviteRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	now time.Time,
	invite PersonalInviteContext,
	status string,
) (err error) {
	invite.Dto.Status = status
	invite.Dto.To.UserID = uid
	inviteUpdates := []dal.Update{
		{Field: "status", Value: status},
		{Field: "to.userID", Value: uid},
	}
	switch status {
	case "active":
		if invite.Dto.Claimed != nil {
			invite.Dto.Claimed = nil
			inviteUpdates = append(inviteUpdates, dal.Update{Field: "claimed", Value: dal.DeleteField})
		}
	case "expired": // Do nothing
	default:
		invite.Dto.Claimed = &now
		inviteUpdates = append(inviteUpdates, dal.Update{Field: "claimed", Value: now})
	}
	if err := invite.Dto.Validate(); err != nil {
		return fmt.Errorf("personal invite record is not valid: %w", err)
	}
	if err = tx.Update(ctx, invite.Key, inviteUpdates); err != nil {
		return err
	}
	return err
}

func updateMemberRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	member dal4contactus.ContactEntry,
	requestMember *briefs4contactus.ContactBase,
	teamMember *briefs4contactus.ContactBase,
) (err error) {
	updates := []dal.Update{
		{Field: "userID", Value: uid},
	}
	updates = updatePersonDetails(&member.Data.ContactBase, requestMember, teamMember, updates)
	if err = tx.Update(ctx, member.Key, updates); err != nil {
		return err
	}
	return err
}

func updateTeamRecord(
	uid, memberID string,
	params *dal4contactus.ContactusTeamWorkerParams,
	requestMember dbmodels.DtoWithID[*briefs4contactus.ContactBase],
) (teamMember *briefs4contactus.ContactBase, err error) {
	if uid == "" {
		panic("required parameter `uid` is empty string")
	}

	inviteToMemberID := memberID[strings.Index(memberID, ":")+1:]
	for contactID, m := range params.TeamModuleEntry.Data.Contacts {
		if contactID == inviteToMemberID {
			m.UserID = uid
			params.TeamModuleEntry.Data.AddUserID(uid)
			params.TeamModuleEntry.Data.AddContact(contactID, m)
			//request.ID.Roles = m.Roles
			//m = request.ID
			m.UserID = uid
			teamMember = &briefs4contactus.ContactBase{
				ContactBrief: *m,
			}
			//team.Members[i] = m
			updatePersonDetails(teamMember, requestMember.Data, teamMember, nil)
			if u, ok := params.TeamModuleEntry.Data.AddUserID(uid); ok {
				params.TeamModuleUpdates = append(params.TeamModuleUpdates, u)
			}
			if m.AddRole(const4contactus.TeamMemberRoleMember) {
				params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{Field: "contacts." + contactID + ".roles", Value: m.Roles})
			}
			break
		}
	}
	if teamMember == nil {
		return teamMember, fmt.Errorf("team member is not found by ID=%s", inviteToMemberID)
	}

	if params.Team.Data.HasUserID(uid) {
		goto UserIdAdded
	}
	params.TeamUpdates = append(params.TeamUpdates, dal.Update{Field: "userIDs", Value: params.Team.Data.UserIDs})
UserIdAdded:
	return teamMember, err
}

func createOrUpdateUserRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	now time.Time,
	user models4userus.UserContext,
	request AcceptPersonalInviteRequest,
	params *dal4contactus.ContactusTeamWorkerParams,
	teamMember *briefs4contactus.ContactBase,
	invite PersonalInviteContext,
) (err error) {
	if teamMember == nil {
		panic("teamMember == nil")
	}
	existingUser := user.Record.Exists()
	if existingUser {
		teamInfo := user.Dto.GetUserTeamInfoByID(request.TeamID)
		if teamInfo != nil {
			return nil
		}
	}

	userTeamInfo := models4userus.UserTeamBrief{
		TeamBrief: params.Team.Data.TeamBrief,
		Roles:     invite.Dto.Roles, // TODO: Validate roles?
	}
	if err = userTeamInfo.Validate(); err != nil {
		return fmt.Errorf("invalid user team info: %w", err)
	}
	user.Dto.Teams[request.TeamID] = &userTeamInfo
	user.Dto.TeamIDs = append(user.Dto.TeamIDs, request.TeamID)
	if existingUser {
		userUpdates := []dal.Update{
			{
				Field: "teams",
				Value: user.Dto.Teams,
			},
		}
		userUpdates = updatePersonDetails(&user.Dto.ContactBase, request.Member.Data, teamMember, userUpdates)
		if err = user.Dto.Validate(); err != nil {
			return fmt.Errorf("user record prepared for update is not valid: %w", err)
		}
		if err = tx.Update(ctx, user.Key, userUpdates); err != nil {
			return fmt.Errorf("failed to update user record: %w", err)
		}
	} else { // New user record
		user.Dto.CreatedAt = now
		user.Dto.Created.Client = request.RemoteClient
		user.Dto.Type = briefs4contactus.ContactTypePerson
		user.Dto.Names = request.Member.Data.Names
		if user.Dto.Names.IsEmpty() && teamMember != nil {
			user.Dto.Names = teamMember.Names
		}
		updatePersonDetails(&user.Dto.ContactBase, request.Member.Data, teamMember, nil)
		if user.Dto.Gender == "" {
			user.Dto.Gender = "unknown"
		}
		if user.Dto.CountryID == "" {
			user.Dto.CountryID = with.UnknownCountryID
		}
		if len(request.Member.Data.Emails) > 0 {
			user.Dto.Emails = request.Member.Data.Emails
		}
		if len(request.Member.Data.Phones) > 0 {
			user.Dto.Phones = request.Member.Data.Phones
		}
		if err = user.Dto.Validate(); err != nil {
			return fmt.Errorf("user record prepared for insert is not valid: %w", err)
		}
		if err = tx.Insert(ctx, user.Record); err != nil {
			return fmt.Errorf("failed to insert user record: %w", err)
		}
	}
	return err
}

func updatePersonDetails(personContact *briefs4contactus.ContactBase, member *briefs4contactus.ContactBase, teamMember *briefs4contactus.ContactBase, updates []dal.Update) []dal.Update {
	if member.Names != nil {
		if personContact.Names == nil {
			personContact.Names = new(person.NameFields)
		}
		if personContact.Names.FirstName == "" {
			name := member.Names.FirstName
			if name == "" {
				name = teamMember.Names.FirstName
			}
			if name != "" {
				personContact.Names.FirstName = name
				if updates != nil {
					updates = append(updates, dal.Update{
						Field: "name.first",
						Value: name,
					})
				}
			}
		}
		if personContact.Names.LastName == "" {
			name := member.Names.LastName
			if name == "" {
				name = teamMember.Names.LastName
			}
			if name != "" {
				personContact.Names.LastName = name
				if updates != nil {
					updates = append(updates, dal.Update{
						Field: "name.last",
						Value: name,
					})
				}
			}
		}
		if personContact.Names.FullName == "" {
			name := member.Names.FullName
			if name == "" {
				name = teamMember.Names.FullName
			}
			if name != "" {
				personContact.Names.FullName = name
				if updates != nil {
					updates = append(updates, dal.Update{
						Field: "name.full",
						Value: name,
					})
				}
			}
		}
	}
	if personContact.Gender == "" || personContact.Gender == "unknown" {
		gender := member.Gender
		if gender == "" || gender == "unknown" {
			gender = teamMember.Gender
		}
		if gender == "" {
			gender = "unknown"
		}
		if personContact.Gender == "" || gender != "unknown" {
			personContact.Gender = member.Gender
			if updates != nil {
				updates = append(updates, dal.Update{
					Field: "gender",
					Value: gender,
				})
			}
		}
	}
	return updates
}
