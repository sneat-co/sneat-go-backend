package facade4retrospectus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dal4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"time"
)

// FixCounts fixes counts
func FixCounts(ctx context.Context, userCtx facade.UserContext, request FixCountsRequest) (err error) {
	uid := userCtx.GetUserID()
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		now := time.Now()
		userRef := dbo4userus.NewUserKey(uid)
		team := dbo4spaceus.NewSpaceEntry(request.SpaceID)
		var retroSpace dal4retrospectus.RetroSpaceEntry
		retroSpace, err = dal4retrospectus.GetRetroSpaceEntry(ctx, tx, request.SpaceID)
		user := new(dbo4userus.UserDbo)
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
		updates := make([]dal.Update, 0, 1)
		if teamInfo == nil {
			if _, ok := retroSpace.Data.UpcomingRetro.ItemsByUserAndType[uid]; ok {
				delete(retroSpace.Data.UpcomingRetro.ItemsByUserAndType, uid)
				if len(retroSpace.Data.UpcomingRetro.ItemsByUserAndType) == 0 {
					retroSpace.Data.UpcomingRetro = nil
					updates = append(updates, dal.Update{Field: "upcomingRetro", Value: dal.DeleteField})
				} else {
					path := fmt.Sprintf("upcomingRetro.itemsByUserAndType.%v", uid)
					updates = append(updates, dal.Update{Field: path, Value: dal.DeleteField})
				}
			}
		} else {
			//for itemType, items := range teamInfo.RetroItems {
			//	count := len(items)
			//	if v, ok := team.Data.UpcomingRetro.ItemsByUserAndType[uid][itemType]; !ok || v != count {
			//		path := fmt.Sprintf("upcomingRetro.itemsByUserAndType.%v.%v", uid, itemType)
			//		if count == 0 {
			//			delete(team.Data.UpcomingRetro.ItemsByUserAndType[uid], itemType)
			//			updates = append(updates, dal.Update{Field: path, Value: dal.DeleteField})
			//		} else {
			//			team.Data.UpcomingRetro.ItemsByUserAndType[uid][itemType] = count
			//			updates = append(updates, dal.Update{Field: path, Value: count})
			//		}
			//	}
			//}
			if len(retroSpace.Data.UpcomingRetro.ItemsByUserAndType[uid]) == 0 {
				delete(retroSpace.Data.UpcomingRetro.ItemsByUserAndType, uid)
				updates = []dal.Update{{Field: "upcomingRetro", Value: dal.DeleteField}}
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
