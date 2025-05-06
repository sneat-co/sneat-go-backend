package dal4calendarium

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

type CalendariumSpaceWorkerParams = dal4spaceus.ModuleSpaceWorkerParams[*dbo4calendarium.CalendariumSpaceDbo]

type HappeningWorkerParams struct {
	*CalendariumSpaceWorkerParams
	Happening        dbo4calendarium.HappeningEntry
	HappeningUpdates []update.Update
}

type HappeningWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *HappeningWorkerParams) (err error)

func RunHappeningSpaceWorker(
	ctx facade.ContextWithUser,
	request dto4calendarium.HappeningRequest,
	happeningWorker HappeningWorker,
) (err error) {
	if err = request.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}
	moduleSpaceWorker := func(
		ctx facade.ContextWithUser,
		tx dal.ReadwriteTransaction,
		moduleSpaceParams *dal4spaceus.ModuleSpaceWorkerParams[*dbo4calendarium.CalendariumSpaceDbo],
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
		if len(params.SpaceModuleUpdates) == 0 &&
			params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring &&
			(len(params.HappeningUpdates) > 0 || params.Happening.Record.HasChanged()) &&
			params.SpaceModuleEntry.Data != nil /* Special case when for example we cancel happening on a specific date */ {
			recurringHappening := params.SpaceModuleEntry.Data.RecurringHappenings[params.Happening.ID]
			if recurringHappening == nil {
				recurringHappening = new(dbo4calendarium.CalendarHappeningBrief)
				if params.SpaceModuleEntry.Data.RecurringHappenings == nil {
					params.SpaceModuleEntry.Data.RecurringHappenings = make(map[string]*dbo4calendarium.CalendarHappeningBrief)
				}
				params.SpaceModuleEntry.Data.RecurringHappenings[params.Happening.ID] = recurringHappening
			}
			recurringHappening.HappeningBase = params.Happening.Data.HappeningBase
			recurringHappening.WithRelated = params.Happening.Data.WithRelated
			moduleSpaceParams.SpaceModuleUpdates = append(moduleSpaceParams.SpaceModuleUpdates,
				update.ByFieldPath([]string{"recurringHappenings", request.HappeningID},
					params.Happening.Data.HappeningBase))
			moduleSpaceParams.SpaceModuleEntry.Record.MarkAsChanged()
		}
		return nil
	}
	return RunCalendariumSpaceWorker(ctx, request.SpaceRequest, moduleSpaceWorker)
}
