package dal4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func RunCalendariumSpaceWorker(
	ctx context.Context,
	user facade.User,
	request dto4spaceus.SpaceRequest,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *CalendariumSpaceWorkerParams) (err error),
) error {
	return dal4spaceus.RunModuleSpaceWorker(ctx, user, request, const4calendarium.ModuleID, new(dbo4calendarium.CalendariumSpaceDbo), worker)
}
