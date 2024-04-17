package facade4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// RemoveTeamMember removes members from a team
func RemoveTeamMember(ctx context.Context, user facade.User, request dto4contactus.ContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	return dal4contactus.RunContactWorker(ctx, user, request,
		func(ctx context.Context, tx dal.ReadwriteTransaction,
			params *dal4contactus.ContactWorkerParams,
		) (err error) {
			return removeTeamMemberTx(ctx, tx, request, params)
		})
}

func removeTeamMemberTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4contactus.ContactRequest,
	params *dal4contactus.ContactWorkerParams,
) (err error) {

	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}

	if params.Contact.Record.Exists() {
		if params.Contact.Data.RemoveRole(const4contactus.TeamMemberRoleMember) {
			params.ContactUpdates = append(params.ContactUpdates, dal.Update{Field: "roles", Value: params.Contact.Data.Roles})
		}
	}

	var memberUserID string
	var membersCount int

	memberUserID, membersCount, err = removeContactBrief(params)
	if err != nil {
		return
	}

	removeMemberFromTeamRecord(params.TeamWorkerParams, memberUserID, membersCount)

	if memberUserID != "" {
		var (
			userRef *dal.Key
		)
		memberUser := models4userus.NewUserContext(memberUserID)
		if err = facade4userus.TxGetUserByID(ctx, tx, memberUser.Record); err != nil {
			return
		}

		update := updateUserRecordOnTeamMemberRemoved(memberUser.Dto, request.TeamID)
		if update != nil {
			if err = txUpdate(ctx, tx, userRef, []dal.Update{*update}); err != nil {
				return err
			}
		}
	}
	return
}

func updateUserRecordOnTeamMemberRemoved(user *models4userus.UserDto, teamID string) *dal.Update {
	delete(user.Teams, teamID)
	user.TeamIDs = slice.RemoveInPlace(teamID, user.TeamIDs)
	return &dal.Update{
		Field: "teams",
		Value: user.Teams,
	}
}

func removeMemberFromTeamRecord(
	params *dal4teamus.TeamWorkerParams,
	contactUserID string,
	membersCount int,
) {
	if contactUserID != "" && slice.Contains(params.Team.Data.UserIDs, contactUserID) {
		params.Team.Data.UserIDs = slice.RemoveInPlace(contactUserID, params.Team.Data.UserIDs)
		params.TeamUpdates = append(params.TeamUpdates, dal.Update{Field: "userIDs", Value: params.Team.Data.UserIDs})
	}
	if params.Team.Data.NumberOf[models4teamus.NumberOfMembersFieldName] != membersCount {
		params.TeamUpdates = append(params.TeamUpdates, params.Team.Data.SetNumberOf(models4teamus.NumberOfMembersFieldName, membersCount))
	}
}

func removeContactBrief(
	params *dal4contactus.ContactWorkerParams,
) (contactUserID string, membersCount int, err error) {

	for id, contactBrief := range params.TeamModuleEntry.Data.Contacts {
		if id == params.Contact.ID {
			params.TeamModuleUpdates = append(params.TeamModuleUpdates, params.TeamModuleEntry.Data.RemoveContact(id))
			if contactBrief.UserID != "" {
				contactUserID = contactBrief.UserID
				userIDs := slice.RemoveInPlace(contactBrief.UserID, params.TeamModuleEntry.Data.UserIDs)
				if len(userIDs) != len(params.TeamModuleEntry.Data.UserIDs) {
					params.TeamModuleEntry.Data.UserIDs = userIDs
					params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{Field: "userIDs", Value: userIDs})
				}
			}
			break
		}
	}
	membersCount = len(params.TeamModuleEntry.Data.GetContactBriefsByRoles(const4contactus.TeamMemberRoleMember))
	return
}
