package facade4retrospectus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"time"
)

// StartRetroReview starts review
func StartRetroReview(ctx facade.ContextWithUser, request RetroRequest) (response RetrospectiveResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = runRetroWorker(ctx, ctx.User(), request,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params facade4meetingus.WorkerParams) error {
			retrospective := params.Meeting.Record.Data().(*dbo4retrospectus.Retrospective)
			retrospective.Stage = dbo4retrospectus.StageReview
			now := time.Now()
			retrospective.TimeLastAction = &now

			var teamRetroUpdates []update.Update
			if teamRetroUpdates, err = moveRetroItemsFromUsers(ctx, tx, params); err != nil {
				return err
			}

			teamRetroUpdates = append(teamRetroUpdates,
				update.ByFieldName("stage", retrospective.Stage),
				update.ByFieldName("timeLastAction", retrospective.TimeLastAction),
			)

			//retrospetiveKey := dal.NewKeyWithID("api4meetingus", ret)
			if err = txUpdateRetrospective(ctx, tx, params.Meeting.Key, retrospective, teamRetroUpdates); err != nil {
				return err
			}

			response.ID = request.MeetingID
			response.Data = retrospective
			return err
		})
	return
}

func moveRetroItemsFromUsers(ctx context.Context, tx dal.ReadwriteTransaction, params facade4meetingus.WorkerParams) (teamRetrosUpdates []update.Update, err error) {
	retrospective := params.Meeting.Record.Data().(*dbo4retrospectus.Retrospective)

	//wg := sync.WaitGroup{}
	userRetroRecords := make([]dal.Record, len(retrospective.Contacts))
	for _, member := range retrospective.Contacts {
		if member.UserID != "" && member.HasRole(const4contactus.SpaceMemberRoleContributor) {
			userRetroRecords = append(userRetroRecords, getUserRetroRecord(member.UserID, params.Space.ID, new(dbo4retrospectus.Retrospective)))
		}
		//}
	}
	if err = tx.GetMulti(ctx, userRetroRecords); err != nil {
		return
	}
	retroItems := make([]interface{}, 0)

	countsByMemberAndType := make(map[string]map[string]int, 0)
	for i, ur := range userRetroRecords {
		if ur == nil {
			continue
		}
		userRetro := ur.Data().(*dbo4retrospectus.Retrospective)

		if len(userRetro.Items) == 0 {
			continue
		}
		userRetroRecord := userRetroRecords[i]
		uid := userRetroRecord.Key().Parent().ID.(string)
		userCounts := make(map[string]int)
		countsByMemberAndType[uid] = userCounts
		for i, retroItem := range userRetro.Items {
			newItemErr := func(message string) error {
				return validation.NewErrBadRecordFieldValue(
					fmt.Sprintf("items[%v]{id=%v}", i, retroItem.ID),
					message,
				)
			}
			if retroItem.Type == "" {
				err = newItemErr("user's retro item has no type")
				return
			}
			if len(retroItem.Children) > 0 {
				err = newItemErr("user's retro item has child items")
				return
			}
			userCounts[retroItem.Type]++
			retroItems = append(retroItems, retroItem)
		}
		if err = txUpdate(ctx, tx, userRetroRecords[i].Key(), []update.Update{
			update.ByFieldName("items", update.DeleteField),
			update.ByFieldName("countsByMemberAndType", update.DeleteField),
		}); err != nil {
			return
		}
	}
	if len(retroItems) > 0 {
		teamRetrosUpdates = []update.Update{
			update.ByFieldName("items", dal.ArrayUnion(retroItems...)),
		}
	}
	return
}
