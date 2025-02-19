package facade4scrumus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-core/facade"
)

func fastSliceRemove(s []string, i int) []string {
	lastIndex := len(s) - 1
	s[i] = s[lastIndex] // Copy last element to index i.
	s[lastIndex] = ""   // Erase last element (write zero value).
	s = s[:lastIndex]   // Truncate slice.
	return s
}

// ThumbUp adds thumb up
func ThumbUp(ctx context.Context, userCtx facade.UserContext, request ThumbUpRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	//err = facade4contactus.RunContactusTeamWorker(ctx, userCtx, request.SpaceRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *facade4contactus.ContactusSpaceWorkerParams) (err error) {
	//	tx.Get
	//	return nil
	//})

	uid := userCtx.GetUserID()

	return runTaskWorker(ctx, userCtx, request.TaskRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params taskWorkerParams) (err error) {
			if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
				return err
			}
			var userContactID string
			for id, member := range params.SpaceModuleEntry.Data.Contacts {
				if member.UserID == uid {
					userContactID = id
				}
			}
			if userContactID == "" {
				return errors.New("not a members of team")
			}

			if params.task == nil {
				return errors.New("task not found by ContactID: " + request.Task)
			}

			if request.Value {
				for _, memberID := range params.task.ThumbUps {
					if memberID == userContactID {
						return nil
					}
				}
				params.task.ThumbUps = append(params.task.ThumbUps, userContactID)
			} else {
				found := false
				for i, memberID := range params.task.ThumbUps {
					if memberID == userContactID {
						found = true
						params.task.ThumbUps = fastSliceRemove(params.task.ThumbUps, i)
						if len(params.task.ThumbUps) == 0 {
							params.task.ThumbUps = nil
						}
						break
					}
				}
				if !found {
					return nil
				}
			}
			return tx.Update(ctx, params.Meeting.Key, []update.Update{
				update.ByFieldName(
					fmt.Sprintf("statuses.%s.byType.%s", request.ContactID, request.Type),
					params.tasks,
				),
			})
		})
}
