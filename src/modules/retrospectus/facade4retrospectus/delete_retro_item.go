package facade4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteRetroItem deletes item from retrospective
func DeleteRetroItem(ctx facade.ContextWithUser, request RetroItemRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	userCtx := ctx.User()
	if request.MeetingID == UpcomingRetrospectiveID {
		return deleteUserRetroItem(ctx, userCtx, request)
	}

	return
}

func deleteUserRetroItem(ctx context.Context, _ facade.UserContext, _ RetroItemRequest) (err error) {
	//uid := userContext.ContactID()
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		//user := new(dbo4userus.User)
		//userKey := facade4userus.NewUserKey(uid)
		//userRecord := dal.NewRecordWithData(userKey, user)
		//now := time.Now()
		//
		//if err = facade4userus.TxGetUserByID(ctx, tx, userRecord); err != nil {
		//	return err
		//}
		//
		//userTeamInfo := user.GetUserSpaceInfoByID(request.Space)
		//if userTeamInfo == nil {
		//	return validation.NewErrBadRequestFieldValue("space", fmt.Sprintf("user does not belong to this team %v, uid=%v", request.Space, uid))
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
		//		if err = updateTeamWithUpcomingRetroUserCounts(ctx, tx, now, uid, request.Space, userTeamInfo.RetroItems); err != nil {
		//			return fmt.Errorf("failed to update team record: %w", err)
		//		}
		//	}
		//	if err = txUpdate(ctx, tx, userKey, []update.Update{
		//		{Field: fmt.Sprintf("????.%v.retroItems.%v", request.Space, request.Role), Value: items},
		//	}); err != nil {
		//		return fmt.Errorf("failed to update user record after deleting retro item: %w", err)
		//	}
		//}
		return err
	})
}
