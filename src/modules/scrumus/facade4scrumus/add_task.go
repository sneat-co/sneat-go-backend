package facade4scrumus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-core/facade"
	"math/rand"
	"strconv"
	"time"
)

var addTaskInTransaction = func(
	ctx context.Context,
	uid string,
	tx dal.ReadwriteTransaction,
	request AddTaskRequest,
	params facade4meetingus.WorkerParams,
) (response *AddTaskResponse, err error) {
	contactusSpace := params.SpaceModuleEntry
	params.Meeting.Record.SetError(nil)
	scrum := params.Meeting.Record.Data().(*dbo4scrumus.Scrum)

	scrumUpdates := make([]dal.Update, 0, 6)

	status := scrum.GetOrCreateStatus(request.ContactID)
	//status.Member.Title = ""
	if contactBrief, ok := contactusSpace.Data.Contacts[request.ContactID]; ok {
		status.Member.ID = request.ContactID
		status.Member.Title = contactBrief.Title
	} else {
		err = fmt.Errorf("unknown contact: %s", request.ContactID)
		return
	}

	var tasks []*dbo4scrumus.Task
	if tasks = status.ByType[request.Type]; tasks == nil {
		tasks = make([]*dbo4scrumus.Task, 0, 1)
	} else {
		// Make sure duplicate calls are discarded
		for _, task := range tasks {
			if task.ID != "" && task.ID == request.Task || task.ID == "" && request.Task == "" && task.Title == request.Title {
				return
			}
		}
	}

	for request.Task == "" {
		randomID := strconv.Itoa(int(rand.Int31n(9999)))
		for _, task := range tasks {
			if task.ID == randomID {
				continue
			}
		}
		request.Task = randomID
	}

	if err = UpdateLastScrumIDIfNeeded(ctx, tx, params); err != nil {
		return nil, err
	}

	tasks = append(tasks, &dbo4scrumus.Task{ID: request.Task, Title: request.Title})
	if params.Meeting.Record.Exists() {
		scrumUpdates = append(scrumUpdates,
			dal.Update{
				Field: "v",
				Value: dal.Increment(1),
			},
			dal.Update{
				Field: fmt.Sprintf("statuses.%s.byType.%s", request.ContactID, request.Type),
				Value: tasks,
			},
			dal.Update{
				Field: fmt.Sprintf("statuses.%s.members", request.ContactID),
				Value: status.Member,
			},
		)
		if request.Type == "risk" {
			scrumUpdates = append(scrumUpdates, dal.Update{
				Field: "risksCount",
				Value: dal.Increment(1),
			})
		}
		if request.Type == "qna" {
			scrumUpdates = append(scrumUpdates, dal.Update{
				Field: "questionsCount",
				Value: dal.Increment(1),
			})
		}
		if err = tx.Update(ctx, params.Meeting.Key, scrumUpdates); err != nil {
			return nil, fmt.Errorf("failed to update scrum record: %v", err)
		}
	} else {
		if request.Type == "risk" {
			scrum.RisksCount = 1
		}
		if status.ByType == nil {
			status.ByType = make(dbo4scrumus.TasksByType, 1)
		}
		status.ByType[request.Type] = tasks
		if err = scrum.Validate(); err != nil {
			return nil, err
		}
		scrumRecord := dal.NewRecordWithData(params.Meeting.Key, scrum)
		if err = tx.Set(ctx, scrumRecord); err != nil {
			return nil, fmt.Errorf("failed to update scrum record: %v", err)
		}
	}

	return &AddTaskResponse{Created: time.Now()}, nil
}

// AddTask adds task
func AddTask(ctx context.Context, userCtx facade.UserContext, request AddTaskRequest) (response *AddTaskResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}

	err = runScrumWorker(ctx, userCtx, request.Request,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params facade4meetingus.WorkerParams) (err error) {
			if err := tx.GetMulti(ctx, []dal.Record{params.SpaceModuleEntry.Record, params.Meeting.Record}); err != nil {
				return err
			}
			response, err = addTaskInTransaction(ctx, params.UserID, tx, request, params)
			return err
		},
	)
	return
}
