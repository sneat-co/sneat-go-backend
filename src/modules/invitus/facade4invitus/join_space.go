package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// JoinSpaceRequest request
type JoinSpaceRequest struct {
	dto4spaceus.SpaceRequest
	InviteID string `json:"inviteID"`
	Pin      string `json:"pin"`
}

// Validate validates request
func (v *JoinSpaceRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if v.InviteID == "" {
		return validation.NewErrRecordIsMissingRequiredField("invite")
	}
	if v.Pin == "" {
		return validation.NewErrRecordIsMissingRequiredField("pin")
	}
	return nil
}

// JoinSpace joins space
func JoinSpace(ctx context.Context, userCtx facade.UserContext, request JoinSpaceRequest) (team *dbo4spaceus.SpaceDbo, err error) {
	if err = request.Validate(); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		return
	}
	uid := userCtx.GetUserID()

	// We intentionally do not use team worker to query both team & user records in parallel
	err = dal4contactus.RunContactusSpaceWorker(ctx, userCtx, request.SpaceRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) error {

		userKey := dbo4userus.NewUserKey(uid)
		userDto := new(dbo4userus.UserDbo)
		userRecord := dal.NewRecordWithData(userKey, userDto)

		inviteKey := NewInviteKey(request.InviteID)
		inviteDto := new(dbo4invitus.InviteDbo)
		inviteRecord := dal.NewRecordWithData(inviteKey, inviteDto)

		if err = params.GetRecords(ctx, tx, userRecord, inviteRecord); err != nil {
			return fmt.Errorf("failed to get some records from DB by ContactID: %w", err)
		}

		if inviteDto.From.UserID == uid {
			err = fmt.Errorf("%w: you can not join using your own invite", facade.ErrForbidden)
			return err
		}

		switch inviteDto.Status {
		case "active": // OK
		case "claimed":
			return fmt.Errorf("%w: the invite already has been claimed", facade.ErrBadRequest)
		case "expired":
			return fmt.Errorf("%w: the invite has expired", facade.ErrBadRequest)
		default:
			return fmt.Errorf("the invite has unknown status: [%s]", inviteDto.Status)
		}

		if inviteDto.Pin == "" {
			return validation.NewErrBadRecordFieldValue("inviteDto.pin", "is empty")
		}

		if inviteDto.Pin != request.Pin {
			return fmt.Errorf("%w: invalid PIN code", facade.ErrForbidden)
		}

		//if team.LastScrum().InviteID != "" {
		//	if err = joinAddUserToLastScrum(ctx, tx, teamKey, *team, uid); err != nil {
		//		return err
		//	}
		//}

		member := dal4contactus.NewContactEntry(inviteDto.SpaceID, inviteDto.To.MemberID)
		if err = tx.Get(ctx, member.Record); err != nil {
			return fmt.Errorf("failed to get member record: %w", err)
		}

		member.Data.UserID = uid
		memberUpdates := []dal.Update{
			{Field: "userID", Value: uid},
		}
		if err = tx.Update(ctx, member.Key, memberUpdates); err != nil {
			return fmt.Errorf("failed to update member record")
		}

		if err = onJoinUpdateMemberBriefInSpaceOrAddIfMissing(
			ctx, tx, params, inviteDto.From.MemberID, member, uid, userDto,
		); err != nil {
			return err
		}
		if err = onJoinAddSpaceToUser(
			ctx, tx, userDto, userRecord, request.SpaceID, team, member,
		); err != nil {
			return fmt.Errorf("failed to update user record: %w", err)
		}
		if err = onJoinUpdateInvite(ctx, tx, uid, inviteKey, inviteDto); err != nil {
			return fmt.Errorf("failed to update invite record: %w", err)
		}
		return nil
	})
	return
}

//func joinAddUserToLastScrum(ctx context.Context, tx dal.ReadwriteTransaction, teamKey *dal.Key, team dbo4spaceus.SpaceDbo, uID string) (err error) {
//	scrumKey := dal.NewKeyWithID("scrums", team.Last.Scrum.ContactID, dal.WithParentKey(teamKey))
//	scrum := new(dbscrum.Scrum)
//	scrumRecord := dal.NewRecordWithData(scrumKey, scrum)
//	if err = tx.Get(ctx, scrumRecord); err != nil {
//		return err
//	}
//	for _, userID := range scrum.UserIDs {
//		if userID == uID {
//			return nil
//		}
//	}
//	scrum.UserIDs = append(scrum.UserIDs, uID)
//	if err = tx.Update(ctx, scrumKey, []dal.Update{{
//		Field: "userIDs",
//		Value: scrum.UserIDs,
//	}}); err != nil {
//		return err
//	}
//	return nil
//}

func onJoinUpdateInvite(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	inviteKey *dal.Key,
	inviteDto *dbo4invitus.InviteDbo,
) (err error) {
	inviteDto.To.UserID = uid
	if err = inviteDto.Validate(); err != nil {
		return fmt.Errorf("invite record is not valid: %w", err)
	}
	inviteUpdates := []dal.Update{
		{Field: "status", Value: "claimed"},
		{Field: "claimed", Value: time.Now()},
		{Field: "toUserID", Value: uid},
	}
	if err = tx.Update(ctx, inviteKey, inviteUpdates); err != nil {
		return fmt.Errorf("failed to update invite record: %w", err)
	}
	return err
}
func onJoinAddSpaceToUser(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userDto *dbo4userus.UserDbo,
	userRecord dal.Record,
	spaceID string,
	space *dbo4spaceus.SpaceDbo,
	member dal4contactus.ContactEntry,
) (err error) {
	var updates []dal.Update
	if userDto == nil {
		panic("required parameter 'userDto' is nil")
	}
	if strings.TrimSpace(spaceID) == "" {
		panic("required parameter 'spaceID' is empty")
	}
	if space == nil {
		panic("required parameter 'space' is nil")
	}
	spaceInfo := userDto.GetUserSpaceInfoByID(spaceID)
	if spaceInfo == nil {
		spaceInfo = &dbo4userus.UserSpaceBrief{
			SpaceBrief: space.SpaceBrief,
			Roles:      member.Data.Roles,
			//MemberType:   "", // TODO: populate?
		}
		userDto.Spaces[spaceID] = spaceInfo
		userDto.SpaceIDs = append(userDto.SpaceIDs, spaceID)
	} else {
		for _, role := range member.Data.Roles {
			hasRole := spaceInfo.HasRole(role)
			if spaceInfo.Title == space.Title && hasRole {
				return // no changes
			}
			spaceInfo.Title = space.Title
			if !hasRole {
				spaceInfo.Roles = append(spaceInfo.Roles, role)
			}
		}
	}
	updates = []dal.Update{
		{
			Field: dbo4spaceus.SpacesFiled,
			Value: userDto.Spaces,
		},
		{
			Field: "spaceIDs",
			Value: userDto.SpaceIDs,
		},
	}
	if len(updates) > 0 {
		if err = userDto.Validate(); err != nil {
			return fmt.Errorf("userDto record is not valid: %w", err)
		}
		if userRecord.Exists() {
			if err = tx.Update(ctx, userRecord.Key(), updates); err != nil {
				return fmt.Errorf("failed to update userDto record: %w", err)
			}
		} else {
			if err = tx.Insert(ctx, userRecord); err != nil {
				return fmt.Errorf("failed to create userDto record: %w", err)
			}
		}
	}
	return
}

func onJoinUpdateMemberBriefInSpaceOrAddIfMissing(
	_ context.Context,
	_ dal.ReadwriteTransaction,
	params *dal4contactus.ContactusSpaceWorkerParams,
	inviterMemberID string,
	member dal4contactus.ContactEntry,
	uid string,
	user *dbo4userus.UserDbo,
) (err error) {
	//var updates []dal.Update
	if strings.TrimSpace(uid) == "" {
		panic("missing required parameter 'uid'")
	}
	if strings.TrimSpace(member.Data.UserID) == "" {
		return validation.NewErrBadRecordFieldValue("userID", "joining member should have populated field 'userID'")
	}
	if member.Data.UserID != uid {
		return validation.NewErrBadRecordFieldValue("userID", fmt.Sprintf("joining member should have same user ContactID as current user, got: {uid=%s, member.Data.UserID=%s}", uid, member.Data.UserID))
	}
	//updates = make([]dal.Update, 0, 2)
	for _, userID := range params.SpaceModuleEntry.Data.UserIDs {
		if userID == uid {
			goto UserIdAddedToUserIDsField
		}
	}

	_, _ = params.Space.Data.AddUserID(uid)
	//if u, ok := params.Space.Data.AddUserID(uid); ok {
	//	updates = append(updates, u)
	//}

UserIdAddedToUserIDsField:

	var memberBrief *briefs4contactus.ContactBrief

	var isValidInviter bool

	for mID, m := range params.SpaceModuleEntry.Data.Contacts {
		if mID == member.ID {
			memberBrief = m
			goto MemberAdded
		} else if m.UserID == uid {
			return fmt.Errorf("current user already joined this team with different contactID=%s", mID)
		}
		if mID == inviterMemberID {
			isValidInviter = true
		}
	}
	if !isValidInviter {
		return fmt.Errorf("supplied inviterMemberID does not belong to the team: %s", inviterMemberID)
	}
	memberBrief = &briefs4contactus.ContactBrief{
		Type:   briefs4contactus.ContactTypePerson,
		Title:  user.Names.GetFullName(),
		Avatar: user.Avatar,
		RolesField: with.RolesField{
			Roles: member.Data.Roles,
		},
		//Emails: user.Emails,
		//Invites: []briefs4memberus.MemberInvite{
		//	{
		//		Channel:         "none",
		//		CreatedBy:       uid,
		//		CreateTime:      time.Now(),
		//		InviterMemberID: inviterMemberID,
		//	},
		//},
	}
	params.SpaceModuleEntry.Data.AddContact(member.ID, memberBrief)
MemberAdded:
	switch memberBrief.UserID {
	case "":
		panic("not implemented")
		//memberBrief.UserID = uid
		//updates = append(updates, dal.Update{
		//	Field: "members",
		//	Value: params.Space.Members,
		//})
	case uid: // Do nothing
	default:
		err = validation.NewErrBadRecordFieldValue("userID", "member already has different userID="+memberBrief.UserID)
	}
	return
}
