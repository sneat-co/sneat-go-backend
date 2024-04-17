package facade4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
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
	return dal4contactus.RunContactusTeamWorker(ctx, user, request.TeamRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusTeamWorkerParams) (err error) {
			return removeTeamMemberTx(ctx, tx, request, params)
		})
}

func removeTeamMemberTx(ctx context.Context, tx dal.ReadwriteTransaction, request dto4contactus.ContactRequest, params *dal4contactus.ContactusTeamWorkerParams) (err error) {

	contact := dal4contactus.NewContactEntry(request.TeamID, request.ContactID)

	if err = params.GetRecords(ctx, tx, contact.Record); err != nil {
		return err
	}

	var memberUserID string

	var contactMatcher = func(contactID string, _ *briefs4contactus.ContactBrief) bool {
		return contactID == request.ContactID
	}

	memberUserID, params.TeamModuleUpdates, err = removeTeamMember(params.Team, params.TeamModuleEntry, contactMatcher)
	if err != nil {
		return
	}

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

func removeTeamMember(
	team dal4teamus.TeamContext,
	contactusTeam dal4contactus.ContactusTeamModuleEntry,
	match func(contactID string, m *briefs4contactus.ContactBrief) bool,
) (memberUserID string, updates []dal.Update, err error) {
	userIds := contactusTeam.Data.UserIDs

	for id, contactBrief := range contactusTeam.Data.Contacts {
		if match(id, contactBrief) {
			if contactBrief.UserID != "" {
				memberUserID = contactBrief.UserID
				userIds = removeTeamUserID(userIds, contactBrief.UserID)
			}

			updates = append(updates, contactusTeam.Data.RemoveContact(id))
		}
	}
	if len(userIds) != len(contactusTeam.Data.UserIDs) {
		contactusTeam.Data.UserIDs = userIds
		if len(userIds) == 0 {
			userIds = nil
		}
		updates = []dal.Update{
			{Field: "userIDs", Value: userIds},
		}
	}
	//updates = append(updates, team.Data.SetNumberOf("contacts", len(contactusTeam.Data.Contacts)))
	updates = append(updates, team.Data.SetNumberOf("members", len(contactusTeam.Data.GetContactBriefsByRoles(const4contactus.TeamMemberRoleMember))))
	return
}

func removeTeamUserID(userIDs []string, userID string) []string {
	uIDs := make([]string, 0, len(userIDs))
	for _, uid := range userIDs {
		if uid != userID {
			uIDs = append(uIDs, uid)
		}
	}
	return uIDs
}
