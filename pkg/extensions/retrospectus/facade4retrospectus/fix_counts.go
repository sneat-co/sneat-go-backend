package facade4retrospectus

import (
	"context"
	"fmt"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	dbo4userus2 "github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/retrospectus/dal4retrospectus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// FixCounts fixes counts
func FixCounts(ctx facade.ContextWithUser, request FixCountsRequest) (err error) {
	userCtx := ctx.User()
	uid := userCtx.GetUserID()
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		now := time.Now()
		userRef := dbo4userus2.NewUserKey(uid)
		team := dbo4spaceus.NewSpaceEntry(request.SpaceID)
		var retroSpace dal4retrospectus.RetroSpaceEntry
		retroSpace, err = dal4retrospectus.GetRetroSpaceEntry(ctx, tx, request.SpaceID)
		user := new(dbo4userus2.UserDbo)
		userRecord := dal.NewRecordWithData(userRef, user)

		if err := tx.GetMulti(ctx, []dal.Record{userRecord, team.Record}); err != nil {
			return err
		}
		if retroSpace.Data.UpcomingRetro == nil {
			retroSpace.Data.UpcomingRetro = &dbo4retrospectus.RetrospectiveCounts{
				ItemsByUserAndType: make(map[string]map[string]int),
			}
		}
		teamInfo := user.GetUserSpaceInfoByID(request.SpaceID)
		updates := make([]update.Update, 0, 1)
		if teamInfo == nil {
			if _, ok := retroSpace.Data.UpcomingRetro.ItemsByUserAndType[uid]; ok {
				delete(retroSpace.Data.UpcomingRetro.ItemsByUserAndType, uid)
				if len(retroSpace.Data.UpcomingRetro.ItemsByUserAndType) == 0 {
					retroSpace.Data.UpcomingRetro = nil
					updates = append(updates, update.ByFieldName("upcomingRetro", update.DeleteField))
				} else {
					path := fmt.Sprintf("upcomingRetro.itemsByUserAndType.%v", uid)
					updates = append(updates, update.ByFieldName(path, update.DeleteField))
				}
			}
		} else {
			//for itemType, items := range teamInfo.RetroItems {
			//	count := len(items)
			//	if v, ok := team.Data.UpcomingRetro.ItemsByUserAndType[uid][itemType]; !ok || v != count {
			//		path := fmt.Sprintf("upcomingRetro.itemsByUserAndType.%v.%v", uid, itemType)
			//		if count == 0 {
			//			delete(team.Data.UpcomingRetro.ItemsByUserAndType[uid], itemType)
			//			updates = append(updates, update.Update{Field: path, Value: update.DeleteField})
			//		} else {
			//			team.Data.UpcomingRetro.ItemsByUserAndType[uid][itemType] = count
			//			updates = append(updates, update.Update{Field: path, Value: count})
			//		}
			//	}
			//}
			if len(retroSpace.Data.UpcomingRetro.ItemsByUserAndType[uid]) == 0 {
				delete(retroSpace.Data.UpcomingRetro.ItemsByUserAndType, uid)
				updates = []update.Update{update.ByFieldName("upcomingRetro", update.DeleteField)}
			}
		}
		if len(updates) > 0 {
			if err = txUpdateSpace(ctx, tx, now, team, updates); err != nil {
				return err
			}
		}
		return nil
	})
}
