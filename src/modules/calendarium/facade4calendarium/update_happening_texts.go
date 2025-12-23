package facade4calendarium

import (
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateHappeningTexts(ctx facade.ContextWithUser, request dto4calendarium.UpdateHappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	if err = dal4calendarium.RunHappeningSpaceWorker(ctx, request.HappeningRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
			return updateHappeningTextsTxWorker(ctx, tx, params, request)
		}); err != nil {
		return fmt.Errorf("failed to update happening: %w", err)
	}
	return nil
}

func updateHappeningTextsTxWorker(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *dal4calendarium.HappeningWorkerParams,
	request dto4calendarium.UpdateHappeningRequest,
) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	if request.Title != params.Happening.Data.Title {
		params.Happening.Data.Title = request.Title
		params.Happening.Record.MarkAsChanged()
		params.HappeningUpdates = append(params.HappeningUpdates,
			update.ByFieldName("title", request.Title))
	}
	if request.Summary != params.Happening.Data.Summary {
		params.Happening.Data.Summary = request.Summary
		params.Happening.Record.MarkAsChanged()
		params.HappeningUpdates = append(params.HappeningUpdates,
			update.ByFieldName("summary", request.Summary))
	}
	if request.Description != params.Happening.Data.Description {
		params.Happening.Data.Description = request.Description
		params.Happening.Record.MarkAsChanged()
		params.HappeningUpdates = append(params.HappeningUpdates,
			update.ByFieldName("description", request.Description))
	}
	if params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
		brief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(request.HappeningID)
		if brief == nil {
			brief = &dbo4calendarium.CalendarHappeningBrief{
				HappeningBase: params.Happening.Data.HappeningBase,
				WithRelated:   params.Happening.Data.WithRelated,
			}
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath([]string{dbo4calendarium.RecurringHappeningsField, request.HappeningID}, brief))
		} else {
			if brief.Title != params.Happening.Data.Title {
				brief.Title = params.Happening.Data.Title
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath([]string{dbo4calendarium.RecurringHappeningsField, request.HappeningID, "title"}, brief.Title))
				params.SpaceModuleEntry.Record.MarkAsChanged()
			}
			if brief.Summary != params.Happening.Data.Summary {
				brief.Summary = params.Happening.Data.Summary
				var summary any
				if brief.Summary == "" {
					summary = update.DeleteField
				} else {
					summary = brief.Summary
				}
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath([]string{dbo4calendarium.RecurringHappeningsField, request.HappeningID, "summary"}, summary))
				params.SpaceModuleEntry.Record.MarkAsChanged()
			}
		}
	}
	return nil
}
