package facade4scrumus

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type taskWorkerParams struct {
	facade4meetingus.WorkerParams
	tasks     dbo4scrumus.Tasks
	task      *dbo4scrumus.Task
	taskIndex int
}

type taskWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params taskWorkerParams) (err error)

func runTaskWorker(ctx context.Context, userCtx facade.UserContext, request TaskRequest, worker taskWorker) (err error) {
	return runScrumWorker(ctx, userCtx, request.Request,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params facade4meetingus.WorkerParams) (err error) {
			taskParams := taskWorkerParams{
				WorkerParams: params,
			}

			scrum := params.Meeting.Record.Data().(*dbo4scrumus.Scrum)

			status := scrum.Statuses[request.ContactID]
			if status == nil {
				return worker(ctx, tx, taskParams)
			}

			if tasks, ok := status.ByType[request.Type]; ok {
				taskParams.tasks = tasks
				for i, task := range tasks {
					if task.ID == request.Task {
						taskParams.task = task
						taskParams.taskIndex = i
						return worker(ctx, tx, taskParams)
					}
				}
			}

			return worker(ctx, tx, taskParams)
		})
}
