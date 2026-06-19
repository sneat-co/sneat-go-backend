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

// DeleteHappening deletes happening
func DeleteHappening(ctx facade.ContextWithUser, request dto4calendarius.HappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	return dal4calendarius.RunHappeningSpaceWorker(ctx, request,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams) (err error) {
			return deleteHappeningTx(ctx, tx, params, request)
		},
	)
}

func deleteHappeningTx(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams, request dto4calendarius.HappeningRequest) (err error) {
	if !params.Happening.Record.Exists() || params.Happening.Data.Type == dbo4calendarius.HappeningTypeRecurring {
		if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
			return
		}
		if happeningBrief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(request.HappeningID); happeningBrief != nil {
			delete(params.SpaceModuleEntry.Data.RecurringHappenings, request.HappeningID)
			params.SpaceModuleUpdates = append(params.SpaceUpdates,
				update.ByFieldPath([]string{dbo4calendarius.RecurringHappeningsField, request.HappeningID}, update.DeleteField))
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
