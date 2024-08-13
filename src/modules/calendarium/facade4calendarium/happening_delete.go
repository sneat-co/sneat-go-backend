package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
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
	switch happening.Data.Type {
	case "":
		return fmt.Errorf("unknown happening type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
	case dbo4calendarium.HappeningTypeSingle:
	case dbo4calendarium.HappeningTypeRecurring:
		happeningBrief := params.SpaceModuleEntry.Data.GetRecurringHappeningBrief(request.HappeningID)

		if happeningBrief != nil {
			delete(params.SpaceModuleEntry.Data.RecurringHappenings, request.HappeningID)
			params.SpaceModuleUpdates = append(params.SpaceUpdates, dal.Update{
				Field: "recurringHappenings." + request.HappeningID,
				Value: dal.DeleteField,
			})
			params.SpaceModuleEntry.Record.MarkAsChanged()
		}
	default:
		return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+happening.Data.Type)
	}
	if happening.Record.Exists() {
		if err = tx.Delete(ctx, happening.Key); err != nil {
			return fmt.Errorf("faield to delete happening record: %w", err)
		}
	}
	return err
}
