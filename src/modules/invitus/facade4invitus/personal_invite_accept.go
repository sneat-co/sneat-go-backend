package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
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
func AcceptPersonalInvite(ctx context.Context, userCtx facade.UserContext, request AcceptPersonalInviteRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	uid := userCtx.GetUserID()

	return dal4contactus.RunContactusSpaceWorker(ctx, userCtx, request.SpaceRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) error {
			var invite PersonalInviteEntry
			var member dal4contactus.ContactEntry

			if invite, member, err = getPersonalInviteRecords(ctx, tx, params, request.InviteID, request.Member.ID); err != nil {
				return err
			}
			if invite.Data.Status != "active" {
				return fmt.Errorf("invite status is not equal to 'active', got '%s'", invite.Data.Status)
			}

			if invite.Data.Pin != request.Pin {
				return fmt.Errorf("%w: pin code does not match", facade.ErrBadRequest)
			}

			user := dbo4userus.NewUserEntry(uid)
			if err = dal4userus.GetUser(ctx, tx, user); err != nil {
				if !dal.IsNotFound(err) {
					return err
				}
			}

			now := time.Now()

			if err = updateInviteRecord(ctx, tx, uid, now, invite, "accepted"); err != nil {
				return fmt.Errorf("failed to update invite record: %w", err)
			}

			var spaceMember *briefs4contactus.ContactBase
			if spaceMember, err = updateSpaceRecord(uid, invite.Data.ToSpaceMemberID, params, request.Member); err != nil {
				return fmt.Errorf("failed to update team record: %w", err)
			}

			memberContext := dal4contactus.NewContactEntry(params.Space.ID, member.ID)

			if err = updateMemberRecord(ctx, tx, uid, memberContext, request.Member.Data, spaceMember); err != nil {
				return fmt.Errorf("failed to update team member record: %w", err)
			}

			if err = createOrUpdateUserRecord(ctx, tx, now, user, request, params, spaceMember, invite); err != nil {
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
	invite PersonalInviteEntry,
	status string,
) (err error) {
	invite.Data.Status = status
	invite.Data.To.UserID = uid
	inviteUpdates := []dal.Update{
		{Field: "status", Value: status},
		{Field: "to.userID", Value: uid},
	}
	switch status {
	case "active":
		if invite.Data.Claimed != nil {
			invite.Data.Claimed = nil
			inviteUpdates = append(inviteUpdates, dal.Update{Field: "claimed", Value: dal.DeleteField})
		}
	case "expired": // Do nothing
	default:
		invite.Data.Claimed = &now
		inviteUpdates = append(inviteUpdates, dal.Update{Field: "claimed", Value: now})
	}
	if err := invite.Data.Validate(); err != nil {
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

func updateSpaceRecord(
	uid, memberID string,
	params *dal4contactus.ContactusSpaceWorkerParams,
	requestMember dbmodels.DtoWithID[*briefs4contactus.ContactBase],
) (teamMember *briefs4contactus.ContactBase, err error) {
	if uid == "" {
		panic("required parameter `uid` is empty string")
	}

	inviteToMemberID := memberID[strings.Index(memberID, ":")+1:]
	for contactID, m := range params.SpaceModuleEntry.Data.Contacts {
		if contactID == inviteToMemberID {
			m.UserID = uid
			params.SpaceModuleEntry.Data.AddUserID(uid)
			params.SpaceModuleEntry.Data.AddContact(contactID, m)
			//request.ContactID.Roles = m.Roles
			//m = request.ContactID
			m.UserID = uid
			teamMember = &briefs4contactus.ContactBase{
				ContactBrief: *m,
			}
			//team.Members[i] = m
			updatePersonDetails(teamMember, requestMember.Data, teamMember, nil)
			if u, ok := params.SpaceModuleEntry.Data.AddUserID(uid); ok {
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, u)
			}
			if m.AddRole(const4contactus.SpaceMemberRoleMember) {
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, dal.Update{Field: "contacts." + contactID + ".roles", Value: m.Roles})
			}
			break
		}
	}
	if teamMember == nil {
		return teamMember, fmt.Errorf("space member is not found by ContactID=%s", inviteToMemberID)
	}

	if params.Space.Data.HasUserID(uid) {
		goto UserIdAdded
	}
	params.SpaceUpdates = append(params.SpaceUpdates, dal.Update{Field: "userIDs", Value: params.Space.Data.UserIDs})
UserIdAdded:
	return teamMember, err
}

func createOrUpdateUserRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	now time.Time,
	user dbo4userus.UserEntry,
	request AcceptPersonalInviteRequest,
	params *dal4contactus.ContactusSpaceWorkerParams,
	teamMember *briefs4contactus.ContactBase,
	invite PersonalInviteEntry,
) (err error) {
	if teamMember == nil {
		panic("spaceMember == nil")
	}
	existingUser := user.Record.Exists()
	if existingUser {
		teamInfo := user.Data.GetUserSpaceInfoByID(request.SpaceID)
		if teamInfo != nil {
			return nil
		}
	}

	userSpaceInfo := dbo4userus.UserSpaceBrief{
		SpaceBrief: params.Space.Data.SpaceBrief,
		Roles:      invite.Data.Roles, // TODO: Validate roles?
	}
	if err = userSpaceInfo.Validate(); err != nil {
		return fmt.Errorf("invalid user team info: %w", err)
	}
	user.Data.Spaces[request.SpaceID] = &userSpaceInfo
	user.Data.SpaceIDs = append(user.Data.SpaceIDs, request.SpaceID)
	if existingUser {
		userUpdates := []dal.Update{
			{
				Field: "spaces",
				Value: user.Data.Spaces,
			},
		}
		userUpdates = updatePersonDetails(&user.Data.ContactBase, request.Member.Data, teamMember, userUpdates)
		if err = user.Data.Validate(); err != nil {
			return fmt.Errorf("user record prepared for update is not valid: %w", err)
		}
		if err = tx.Update(ctx, user.Key, userUpdates); err != nil {
			return fmt.Errorf("failed to update user record: %w", err)
		}
	} else { // New user record
		user.Data.CreatedAt = now
		user.Data.Created.Client = request.RemoteClient
		user.Data.Type = briefs4contactus.ContactTypePerson
		user.Data.Names = request.Member.Data.Names
		if user.Data.Names.IsEmpty() {
			user.Data.Names = teamMember.Names
		}
		updatePersonDetails(&user.Data.ContactBase, request.Member.Data, teamMember, nil)
		if user.Data.Gender == "" {
			user.Data.Gender = "unknown"
		}
		if user.Data.CountryID == "" {
			user.Data.CountryID = with.UnknownCountryID
		}
		if len(request.Member.Data.Emails) > 0 {
			user.Data.Emails = request.Member.Data.Emails
		}
		if len(request.Member.Data.Phones) > 0 {
			user.Data.Phones = request.Member.Data.Phones
		}
		if err = user.Data.Validate(); err != nil {
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
