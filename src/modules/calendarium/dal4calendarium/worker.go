package dal4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type CalendariumTeamWorkerParams = dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo]

type HappeningWorkerParams struct {
	*CalendariumTeamWorkerParams
	Happening        dbo4calendarium.HappeningEntry
	HappeningUpdates []dal.Update
}

type HappeningWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *HappeningWorkerParams) (err error)

func RunHappeningTeamWorker(
	ctx context.Context,
	user facade.User,
	request dto4calendarium.HappeningRequest,
	happeningWorker HappeningWorker,
) (err error) {
	moduleTeamWorker := func(
		ctx context.Context,
		tx dal.ReadwriteTransaction,
		moduleTeamParams *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo],
	) (err error) {
		params := &HappeningWorkerParams{
			CalendariumTeamWorkerParams: moduleTeamParams,
			Happening:                   dbo4calendarium.NewHappeningEntry(request.TeamID, request.HappeningID),
		}
		if err = tx.Get(ctx, params.Happening.Record); err != nil {
			if dal.IsNotFound(err) {
				params.Happening.Data.Type = request.HappeningType
			} else {
				return fmt.Errorf("failed to get happening: %w", err)
			}
		}

		if err = happeningWorker(ctx, tx, params); err != nil {
			return err
		}
		if len(params.HappeningUpdates) > 0 {
			if err = tx.Update(ctx, params.Happening.Key, params.HappeningUpdates); err != nil {
				return fmt.Errorf("failed to update happening record: %w", err)
			}
		}
		if len(params.TeamModuleUpdates) == 0 && params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring && (len(params.HappeningUpdates) > 0 || params.Happening.Record.HasChanged()) {
			recurringHappening := params.TeamModuleEntry.Data.RecurringHappenings[params.Happening.ID]
			recurringHappening.HappeningBrief = params.Happening.Data.HappeningBrief
			recurringHappening.WithRelated = params.Happening.Data.WithRelated
			moduleTeamParams.TeamModuleUpdates = append(moduleTeamParams.TeamModuleUpdates, dal.Update{
				Field: "recurringHappenings." + request.HappeningID,
				Value: params.Happening.Data.HappeningBrief,
			})
		}
		return nil
	}
	return RunCalendariumTeamWorker(ctx, user, request.TeamRequest, moduleTeamWorker)
}
