package delays4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
)

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Delayer) {
	updateHappeningBriefDelayer = mustRegisterFunc(delayUpdateHappeningBriefName, delayedUpdateHappeningBrief)
}

func DelayUpdateHappeningBrief(ctx context.Context, spaceID coretypes.SpaceID, happeningID string) (err error) {
	return updateHappeningBriefDelayer.EnqueueWork(ctx, delaying.With(const4calendarium.QueueHappeningBrief, delayUpdateHappeningBriefName, 0), spaceID, happeningID)
}

const delayUpdateHappeningBriefName = "updateHappeningBriefDelayer"

var updateHappeningBriefDelayer delaying.Delayer

func delayedUpdateHappeningBrief(ctx context.Context, spaceID coretypes.SpaceID, happeningID string) (err error) {
	request := dto4spaceus.SpaceItemRequest{
		SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: spaceID},
		ID:           happeningID,
	}
	return dal4spaceus.RunSpaceItemWorker(ctx,
		facade.NewUserContext(""),
		request,
		const4calendarium.ModuleID,
		new(dbo4calendarium.CalendariumSpaceDbo),
		const4calendarium.HappeningsCollection,
		new(dbo4calendarium.HappeningDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceItemWorkerParams[*dbo4calendarium.CalendariumSpaceDbo, *dbo4calendarium.HappeningDbo]) (err error) {
			if err = params.GetRecords(ctx, tx, params.SpaceModuleEntry.Record); err != nil {
				return err
			}
			if params.SpaceItem.Record.Exists() {
				brief := dbo4calendarium.CalendarHappeningBrief{
					HappeningBrief: params.SpaceItem.Data.HappeningBrief,
					WithRelated:    params.SpaceItem.Data.WithRelated,
				}
				params.SpaceModuleEntry.Data.RecurringHappenings[params.SpaceItem.ID] = &brief
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath([]string{"recurringHappenings", params.SpaceItem.ID}, brief))
			} else if params.SpaceModuleEntry.Data.RecurringHappenings[params.SpaceItem.ID] != nil {
				delete(params.SpaceModuleEntry.Data.RecurringHappenings, params.SpaceItem.ID)
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.DeleteByFieldPath("recurringHappenings", happeningID))
			}
			return nil
		},
	)
}
