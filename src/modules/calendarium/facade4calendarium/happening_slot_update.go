package facade4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateSlot(ctx context.Context, user facade.User, request dto4calendarium.HappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) (err error) {
		if err = params.GetRecords(ctx, tx); err != nil {
			return err
		}
		if params.Happening.Data.Type == dbo4calendarium.HappeningTypeSingle {
			params.Happening.Data.Slots[0] = &request.Slot
			params.HappeningUpdates = []dal.Update{
				{
					Field: "slots",
					Value: params.Happening.Data.Slots,
				},
			}
		}
		return
	}

	return dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, worker)
}
