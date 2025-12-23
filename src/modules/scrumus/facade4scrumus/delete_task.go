package facade4scrumus

import (
	"context"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteTask deletes task
func DeleteTask(ctx facade.ContextWithUser, request DeleteTaskRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	return runTaskWorker(ctx, ctx.User(), request,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params taskWorkerParams) (err error) {
			if params.task == nil {
				//err = errors.New("task not found by ContactID: " + request.Task)
				return
			}
			tasks := make([]*dbo4scrumus.Task, 0, len(params.tasks))
			for _, task := range params.tasks {
				if task.ID != request.Task {
					tasks = append(tasks, task)
				}
			}
			if len(tasks) == len(params.tasks) {
				return nil
			}
			var updateValue interface{}
			if len(tasks) == 0 {
				updateValue = update.DeleteField
			} else {
				if err = dbo4scrumus.ValidateTasks(tasks); err != nil {
					return err
				}
				updateValue = tasks
			}
			updates := []update.Update{
				update.ByFieldName("v", dal.Increment(1)),
				update.ByFieldName(fmt.Sprintf("statuses.%s.byType.%s", request.ContactID, request.Type), updateValue),
			}
			if request.Type == "risk" {
				updates = append(updates, update.ByFieldName("risksCount", dal.Increment(-1)))
			}
			if request.Type == "qna" {
				updates = append(updates, update.ByFieldName("questionsCount", dal.Increment(-1)))
			}
			return tx.Update(ctx, params.Meeting.Key, updates)
		})
}
