package facade4calendarius

import (
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dal4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateHappeningTexts(ctx facade.ContextWithUser, request dto4calendarius.UpdateHappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	if err = dal4calendarius.RunHappeningSpaceWorker(ctx, request.HappeningRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams) error {
			return updateHappeningTextsTxWorker(ctx, tx, params, request)
		}); err != nil {
		return fmt.Errorf("failed to update happening: %w", err)
	}
	return nil
}

func updateHappeningTextsTxWorker(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *dal4calendarius.HappeningWorkerParams,
	request dto4calendarius.UpdateHappeningRequest,
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
	if params.Happening.Data.Type == dbo4calendarius.HappeningTypeRecurring {
		brief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(request.HappeningID)
		if brief == nil {
			brief = &dbo4calendarius.CalendarHappeningBrief{
				HappeningBase: params.Happening.Data.HappeningBase,
				WithRelated:   params.Happening.Data.WithRelated,
			}
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath([]string{dbo4calendarius.RecurringHappeningsField, request.HappeningID}, brief))
		} else {
			if brief.Title != params.Happening.Data.Title {
				brief.Title = params.Happening.Data.Title
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath([]string{dbo4calendarius.RecurringHappeningsField, request.HappeningID, "title"}, brief.Title))
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
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath([]string{dbo4calendarius.RecurringHappeningsField, request.HappeningID, "summary"}, summary))
				params.SpaceModuleEntry.Record.MarkAsChanged()
			}
		}
	}
	return nil
}
