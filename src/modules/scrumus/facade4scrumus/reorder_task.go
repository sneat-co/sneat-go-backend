package facade4scrumus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// ReorderTask reorders tasks
func ReorderTask(ctx context.Context, userContext facade.User, request ReorderTaskRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	uid := userContext.GetID()
	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		var params facade4meetingus.WorkerParams
		scrum := dbo4scrumus.Scrum{}
		if params, err = facade4meetingus.GetMeetingAndSpace(ctx, tx, uid, request.SpaceID, request.MeetingID, MeetingRecordFactory{}); err != nil {
			return
		}
		if !params.Meeting.Record.Exists() {
			err = errors.New("scrum record not found by ContactID: " + request.MeetingID)
			return
		}
		status := scrum.Statuses[request.ContactID]
		if status == nil {
			err = errors.New("status not found by members ContactID: " + request.ContactID)
			return
		}
		tasks := status.ByType[request.Type]
		if len(tasks) <= request.From {
			return fmt.Errorf("len(tasks) <= request.From: %d < %d", len(tasks), request.From)
		}
		if len(tasks) <= request.To {
			return fmt.Errorf("len(tasks) <= request.To: %d < %d", len(tasks), request.To)
		}
		task := tasks[request.From]
		if task.ID == request.Task && len(tasks) == request.Len {
			if request.To > request.From {
				for i := request.From; i < request.To; i++ {
					tasks[i] = tasks[i+1]
				}
				tasks[request.To] = task
			} else if request.To < request.From {
				for i := request.From; i > request.To; i-- {
					tasks[i] = tasks[i-1]
				}
				tasks[request.To] = task
			}
		} else {
			return errors.New("reordering on already changed list is not implemented yet")
		}

		return tx.Update(ctx, params.Meeting.Key, []dal.Update{
			{
				Field: "v",
				Value: dal.Increment(1),
			},
			{
				Field: fmt.Sprintf("statuses.%s.byType.%s", request.ContactID, request.Type),
				Value: tasks,
			},
		})
	})
}
