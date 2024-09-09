package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteHappening deletes happening
func DeleteHappening(ctx context.Context, userCtx facade.UserContext, request dto4calendarium.HappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	worker := func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
		return deleteHappeningTx(ctx, tx, request, params)
	}

	return dal4calendarium.RunHappeningSpaceWorker(ctx, userCtx, request, worker)
}

func deleteHappeningTx(ctx context.Context, tx dal.ReadwriteTransaction, request dto4calendarium.HappeningRequest, params *dal4calendarium.HappeningWorkerParams) (err error) {
	happening := params.Happening

	if !happening.Record.Exists() || happening.Data.Type == dbo4calendarium.HappeningTypeRecurring {
		if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
			return
		}
		if happeningBrief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(request.HappeningID); happeningBrief != nil {
			delete(params.SpaceModuleEntry.Data.RecurringHappenings, request.HappeningID)
			params.SpaceModuleUpdates = append(params.SpaceUpdates, dal.Update{
				Field: "recurringHappenings." + request.HappeningID,
				Value: dal.DeleteField,
			})
			params.SpaceModuleEntry.Record.MarkAsChanged()
		}
	}

	if happening.Record.Exists() {
		if err = tx.Delete(ctx, happening.Key); err != nil {
			return fmt.Errorf("faield to delete happening record: %w", err)
		}
	}
	return err
}
