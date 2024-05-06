package dal4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type CalendariumTeamWorkerParams = dal4teamus.ModuleTeamWorkerParams[*models4calendarium.CalendariumTeamDbo]

type HappeningWorkerParams struct {
	CalendariumTeamWorkerParams
	Happening        models4calendarium.HappeningContext
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
		moduleTeamParams *dal4teamus.ModuleTeamWorkerParams[*models4calendarium.CalendariumTeamDbo],
	) (err error) {
		params := &HappeningWorkerParams{
			CalendariumTeamWorkerParams: *moduleTeamParams,
			Happening:                   models4calendarium.NewHappeningContext(request.TeamID, request.HappeningID),
		}
		if err = tx.Get(ctx, params.Happening.Record); err != nil {
			if dal.IsNotFound(err) {
				params.Happening.Dbo.Type = request.HappeningType
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
		if len(params.TeamModuleUpdates) == 0 && params.Happening.Dbo.Type == models4calendarium.HappeningTypeRecurring && (len(params.HappeningUpdates) > 0 || params.Happening.Record.HasChanged()) {
			recurringHappening := params.TeamModuleEntry.Data.RecurringHappenings[params.Happening.ID]
			recurringHappening.HappeningBrief = params.Happening.Dbo.HappeningBrief
			recurringHappening.WithRelated = params.Happening.Dbo.WithRelated
			moduleTeamParams.TeamModuleUpdates = append(moduleTeamParams.TeamModuleUpdates, dal.Update{
				Field: "recurringHappenings." + request.HappeningID,
				Value: params.Happening.Dbo.HappeningBrief,
			})
		}
		return nil
	}
	return RunCalendariumTeamWorker(ctx, user, request.TeamRequest, moduleTeamWorker)
}
