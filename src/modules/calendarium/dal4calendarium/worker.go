package dal4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

type CalendariumSpaceWorkerParams = dal4teamus.ModuleSpaceWorkerParams[*dbo4calendarium.CalendariumSpaceDbo]

type HappeningWorkerParams struct {
	*CalendariumSpaceWorkerParams
	Happening        dbo4calendarium.HappeningEntry
	HappeningUpdates []dal.Update
}

type HappeningWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *HappeningWorkerParams) (err error)

func RunHappeningSpaceWorker(
	ctx context.Context,
	user facade.User,
	request dto4calendarium.HappeningRequest,
	happeningWorker HappeningWorker,
) (err error) {
	if err = request.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}
	moduleSpaceWorker := func(
		ctx context.Context,
		tx dal.ReadwriteTransaction,
		moduleSpaceParams *dal4teamus.ModuleSpaceWorkerParams[*dbo4calendarium.CalendariumSpaceDbo],
	) (err error) {
		params := &HappeningWorkerParams{
			CalendariumSpaceWorkerParams: moduleSpaceParams,
			Happening:                    dbo4calendarium.NewHappeningEntry(request.SpaceID, request.HappeningID),
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
		if len(params.SpaceModuleUpdates) == 0 && params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring && (len(params.HappeningUpdates) > 0 || params.Happening.Record.HasChanged()) {
			recurringHappening := params.SpaceModuleEntry.Data.RecurringHappenings[params.Happening.ID]
			recurringHappening.HappeningBrief = params.Happening.Data.HappeningBrief
			recurringHappening.WithRelated = params.Happening.Data.WithRelated
			moduleSpaceParams.SpaceModuleUpdates = append(moduleSpaceParams.SpaceModuleUpdates, dal.Update{
				Field: "recurringHappenings." + request.HappeningID,
				Value: params.Happening.Data.HappeningBrief,
			})
		}
		return nil
	}
	return RunCalendariumSpaceWorker(ctx, user, request.SpaceRequest, moduleSpaceWorker)
}
