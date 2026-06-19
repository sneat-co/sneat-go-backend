package dal4calendarius

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/const4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-core/facade"
)

func RunCalendariusSpaceWorker(
	ctx facade.ContextWithUser,
	request dto4spaceus.SpaceRequest,
	worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *CalendariusSpaceWorkerParams) (err error),
) error {
	return dal4spaceus.RunModuleSpaceWorkerWithUserCtx(ctx,
		request.SpaceID, const4calendarius.ExtensionID,
		new(dbo4calendarius.CalendariusSpaceDbo),
		worker)
}
