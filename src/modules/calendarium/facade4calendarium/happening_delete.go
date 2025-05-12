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

// DeleteHappening deletes happening
func DeleteHappening(ctx facade.ContextWithUser, request dto4calendarium.HappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	return dal4calendarium.RunHappeningSpaceWorker(ctx, request,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
			return deleteHappeningTx(ctx, tx, params, request)
		},
	)
}

func deleteHappeningTx(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.HappeningRequest) (err error) {
	if !params.Happening.Record.Exists() || params.Happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
		if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
			return
		}
		if happeningBrief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(request.HappeningID); happeningBrief != nil {
			delete(params.SpaceModuleEntry.Data.RecurringHappenings, request.HappeningID)
			params.SpaceModuleUpdates = append(params.SpaceUpdates,
				update.ByFieldPath([]string{dbo4calendarium.RecurringHappeningsField, request.HappeningID}, update.DeleteField))
			params.SpaceModuleEntry.Record.MarkAsChanged()
		}
	}
	if params.Happening.Record.Exists() {
		if err = tx.Delete(ctx, params.Happening.Key); err != nil {
			return fmt.Errorf("faield to delete happening record: %w", err)
		}
	}
	return err
}
