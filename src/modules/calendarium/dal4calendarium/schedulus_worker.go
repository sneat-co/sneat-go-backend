package dal4calendarium

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
)

func RunCalendariumSpaceWorker(
	ctx facade.ContextWithUser,
	request dto4spaceus.SpaceRequest,
	worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *CalendariumSpaceWorkerParams) (err error),
) error {
	return dal4spaceus.RunModuleSpaceWorkerWithUserCtx(ctx,
		request.SpaceID, const4calendarium.ExtensionID,
		new(dbo4calendarium.CalendariumSpaceDbo),
		worker)
}
