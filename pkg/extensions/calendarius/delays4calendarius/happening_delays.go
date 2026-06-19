package delays4calendarius

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/const4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
)

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Delayer) {
	updateHappeningBriefDelayer = mustRegisterFunc(delayUpdateHappeningBriefName, delayedUpdateHappeningBrief)
}

func DelayUpdateHappeningBrief(ctx context.Context, userID string, spaceID coretypes.SpaceID, happeningID string) (err error) {
	return updateHappeningBriefDelayer.EnqueueWork(ctx, delaying.With(const4calendarius.QueueHappeningBrief, delayUpdateHappeningBriefName, 0), userID, spaceID, happeningID)
}

const delayUpdateHappeningBriefName = "updateHappeningBriefDelayer"

var updateHappeningBriefDelayer delaying.Delayer

func delayedUpdateHappeningBrief(ctx context.Context, userID string, spaceID coretypes.SpaceID, happeningID string) (err error) {
	request := dto4spaceus.SpaceItemRequest{
		SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: spaceID},
		ID:           happeningID,
	}
	ctxWithUser := facade.NewContextWithUserID(ctx, userID)
	return dal4spaceus.RunSpaceItemWorker(ctxWithUser,
		request,
		const4calendarius.ExtensionID,
		new(dbo4calendarius.CalendariusSpaceDbo),
		const4calendarius.HappeningsCollection,
		new(dbo4calendarius.HappeningDbo),
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceItemWorkerParams[*dbo4calendarius.CalendariusSpaceDbo, *dbo4calendarius.HappeningDbo]) (err error) {
			if err = params.GetRecords(ctx, tx, params.SpaceModuleEntry.Record); err != nil {
				return err
			}
			if params.SpaceItem.Record.Exists() {
				brief := dbo4calendarius.CalendarHappeningBrief{
					HappeningBase: params.SpaceItem.Data.HappeningBase,
					WithRelated:   params.SpaceItem.Data.WithRelated,
				}
				params.SpaceModuleEntry.Data.RecurringHappenings[params.SpaceItem.ID] = &brief
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldPath([]string{dbo4calendarius.RecurringHappeningsField, params.SpaceItem.ID}, brief))
			} else if params.SpaceModuleEntry.Data.RecurringHappenings[params.SpaceItem.ID] != nil {
				delete(params.SpaceModuleEntry.Data.RecurringHappenings, params.SpaceItem.ID)
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.DeleteByFieldPath(dbo4calendarius.RecurringHappeningsField, happeningID))
			}
			return nil
		},
	)
}
