package facade4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteRetroItem deletes item from retrospective
func DeleteRetroItem(ctx context.Context, userContext facade.User, request RetroItemRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	if request.MeetingID == UpcomingRetrospectiveID {
		return deleteUserRetroItem(ctx, userContext, request)
	}

	return
}

func deleteUserRetroItem(ctx context.Context, userContext facade.User, request RetroItemRequest) (err error) {
	//uid := userContext.ContactID()
	db := facade.GetDatabase(ctx)
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		//user := new(models4userus.User)
		//userKey := facade4userus.NewUserKey(uid)
		//userRecord := dal.NewRecordWithData(userKey, user)
		//now := time.Now()
		//
		//if err = facade4userus.TxGetUserByID(ctx, tx, userRecord); err != nil {
		//	return err
		//}
		//
		//userTeamInfo := user.GetUserTeamInfoByID(request.TeamID)
		//if userTeamInfo == nil {
		//	return validation.NewErrBadRequestFieldValue("team", fmt.Sprintf("user does not belong to this team %v, uid=%v", request.TeamID, uid))
		//}
		//
		//if userTeamInfo.RetroItems == nil {
		//	return nil
		//}
		//
		//existingItems := userTeamInfo.RetroItems[request.Role]
		//if len(existingItems) == 0 {
		//	return nil
		//}
		//
		//items := make([]*dbretro.RetroItem, 0, len(existingItems))
		//for _, item := range existingItems {
		//	if item.ContactID != request.Item {
		//		items = append(items, item)
		//	}
		//}
		//if len(items) != len(existingItems) {
		//	userTeamInfo.RetroItems[request.Role] = items
		//	if request.MeetingID == "upcoming" {
		//		if err = updateTeamWithUpcomingRetroUserCounts(ctx, tx, now, uid, request.TeamID, userTeamInfo.RetroItems); err != nil {
		//			return fmt.Errorf("failed to update team record: %w", err)
		//		}
		//	}
		//	if err = txUpdate(ctx, tx, userKey, []dal.Update{
		//		{Field: fmt.Sprintf("????.%v.retroItems.%v", request.TeamID, request.Role), Value: items},
		//	}); err != nil {
		//		return fmt.Errorf("failed to update user record after deleting retro item: %w", err)
		//	}
		//}
		return err
	})
	return
}
