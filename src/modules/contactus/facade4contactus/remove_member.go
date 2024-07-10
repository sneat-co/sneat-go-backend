package facade4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// RemoveSpaceMember removes members from a team
func RemoveSpaceMember(ctx context.Context, user facade.User, request dto4contactus.ContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	return dal4contactus.RunContactWorker(ctx, user, request,
		func(ctx context.Context, tx dal.ReadwriteTransaction,
			params *dal4contactus.ContactWorkerParams,
		) (err error) {
			return removeSpaceMemberTx(ctx, tx, request, params)
		})
}

func removeSpaceMemberTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4contactus.ContactRequest,
	params *dal4contactus.ContactWorkerParams,
) (err error) {

	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}

	if params.Contact.Record.Exists() {
		if params.Contact.Data.RemoveRole(const4contactus.SpaceMemberRoleMember) {
			params.ContactUpdates = append(params.ContactUpdates, dal.Update{Field: "roles", Value: params.Contact.Data.Roles})
		}
	}

	var memberUserID string
	var membersCount int

	memberUserID, membersCount, err = removeContactBrief(params)
	if err != nil {
		return
	}

	removeMemberFromSpaceRecord(params.SpaceWorkerParams, memberUserID, membersCount)

	if memberUserID != "" {
		var (
			userRef *dal.Key
		)
		memberUser := dbo4userus.NewUserEntry(memberUserID)
		if err = facade4userus.TxGetUserByID(ctx, tx, memberUser.Record); err != nil {
			return
		}

		update := updateUserRecordOnSpaceMemberRemoved(memberUser.Data, request.SpaceID)
		if update != nil {
			if err = txUpdate(ctx, tx, userRef, []dal.Update{*update}); err != nil {
				return err
			}
		}
	}
	return
}

func updateUserRecordOnSpaceMemberRemoved(user *dbo4userus.UserDbo, teamID string) *dal.Update {
	delete(user.Spaces, teamID)
	user.SpaceIDs = slice.RemoveInPlace(teamID, user.SpaceIDs)
	return &dal.Update{
		Field: "spaces",
		Value: user.Spaces,
	}
}

func removeMemberFromSpaceRecord(
	params *dal4teamus.SpaceWorkerParams,
	contactUserID string,
	membersCount int,
) {
	if contactUserID != "" && slice.Contains(params.Space.Data.UserIDs, contactUserID) {
		params.Space.Data.UserIDs = slice.RemoveInPlace(contactUserID, params.Space.Data.UserIDs)
		params.SpaceUpdates = append(params.SpaceUpdates, dal.Update{Field: "userIDs", Value: params.Space.Data.UserIDs})
	}
	//if params.Space.Data.NumberOf[dbo4teamus.NumberOfMembersFieldName] != membersCount {
	//	params.SpaceUpdates = append(params.SpaceUpdates, params.Space.Data.SetNumberOf(dbo4teamus.NumberOfMembersFieldName, membersCount))
	//}
}

func removeContactBrief(
	params *dal4contactus.ContactWorkerParams,
) (contactUserID string, membersCount int, err error) {

	for id, contactBrief := range params.SpaceModuleEntry.Data.Contacts {
		if id == params.Contact.ID {
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, params.SpaceModuleEntry.Data.RemoveContact(id))
			if contactBrief.UserID != "" {
				contactUserID = contactBrief.UserID
				userIDs := slice.RemoveInPlace(contactBrief.UserID, params.SpaceModuleEntry.Data.UserIDs)
				if len(userIDs) != len(params.SpaceModuleEntry.Data.UserIDs) {
					params.SpaceModuleEntry.Data.UserIDs = userIDs
					params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, dal.Update{Field: "userIDs", Value: userIDs})
				}
			}
			break
		}
	}
	membersCount = len(params.SpaceModuleEntry.Data.GetContactBriefsByRoles(const4contactus.SpaceMemberRoleMember))
	return
}
