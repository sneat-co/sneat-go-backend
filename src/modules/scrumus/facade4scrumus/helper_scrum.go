package facade4scrumus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dal4scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params facade4meetingus.WorkerParams) (err error)

func runScrumWorker(ctx context.Context, userCtx facade.UserContext, request facade4meetingus.Request, worker worker) error {
	return facade4meetingus.RunMeetingWorker(ctx, userCtx.GetUserID(), request, MeetingRecordFactory{}, worker)
}

// UpdateLastScrumIDIfNeeded updates scrum if needed
func UpdateLastScrumIDIfNeeded(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params facade4meetingus.WorkerParams,
) (err error) {

	scrumSpaceUpdates := make([]dal.Update, 0, 1)
	scrumID := params.Meeting.GetID()
	scrum := params.Meeting.Record.Data().(*dbo4scrumus.Scrum)

	var scrumSpace dal4scrumus.ScrumSpaceEntry
	scrumSpace, err = dal4scrumus.GetScrumSpace(ctx, tx, params.Space.ID)
	if err != nil && !dal.IsNotFound(err) {
		return
	}
	if lastScrum := scrumSpace.Data.Last; lastScrum != nil && lastScrum.ID != "" && lastScrum.ID < scrumID {
		if scrum.ScrumIDs == nil {
			scrum.ScrumIDs = &dbo4scrumus.ScrumIDs{}
		}
		scrum.ScrumIDs.Prev = lastScrum.ID
		prevScrumKey := dal.NewKeyWithParentAndID(params.Space.Key, "api4meetingus", lastScrum.ID)
		prevScrum := new(dbo4scrumus.Scrum)
		prevScrumRecord := dal.NewRecordWithData(prevScrumKey, prevScrum)
		if err = tx.Get(ctx, prevScrumRecord); err != nil {
			return
		}
		if prevScrum.ScrumIDs == nil || prevScrum.ScrumIDs.Next == "" {
			if err = prevScrum.Validate(); err != nil {
				return
			}
			prevScrumUpdates := []dal.Update{{Field: "scrumIds.next", Value: scrumID}}
			if err = tx.Update(ctx, prevScrumKey, prevScrumUpdates); err != nil {
				return
			}
		}
	}
	if scrumSpace.Data.Last == nil || scrumSpace.Data.Last.ID < scrumID {
		scrumSpace.Data.Last = &dbo4spaceus.SpaceMeetingInfo{
			ID:       scrumID,
			Stage:    "planning",
			Started:  scrum.Started,
			Finished: scrum.Finished,
		}
		scrumSpaceUpdates = append(scrumSpaceUpdates, dal.Update{
			Field: "last",
			Value: scrumSpace.Data.Last,
		})
	}
	if len(scrumSpaceUpdates) > 0 {
		if err = scrumSpace.Data.Validate(); err != nil {
			return
		}
		if scrumSpace.Record.Exists() {
			if err = tx.Update(ctx, scrumSpace.Key, scrumSpaceUpdates); err != nil {
				return fmt.Errorf("failed to update scrum team record: %w", err)
			}
		} else if err = tx.Insert(ctx, scrumSpace.Record); err != nil {
			return fmt.Errorf("failed to insert new scrum team record")
		}
	}
	return
}
