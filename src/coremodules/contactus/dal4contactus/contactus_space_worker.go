package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ContactusSpaceWorkerParams = dal4spaceus.ModuleSpaceWorkerParams[*dbo4contactus.ContactusSpaceDbo]

func NewContactusSpaceWorkerParams(userCtx facade.UserContext, spaceID string) *ContactusSpaceWorkerParams {
	teamWorkerParams := dal4spaceus.NewSpaceWorkerParams(userCtx, spaceID)
	return dal4spaceus.NewSpaceModuleWorkerParams(const4contactus.ModuleID, teamWorkerParams, new(dbo4contactus.ContactusSpaceDbo))
}

func RunReadonlyContactusSpaceWorker(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, params *ContactusSpaceWorkerParams) (err error),
) error {
	return dal4spaceus.RunReadonlyModuleSpaceWorker(ctx, userCtx, request, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo), worker)
}

type ContactusModuleWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactusSpaceWorkerParams) (err error)

func RunContactusSpaceWorker(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	worker ContactusModuleWorker,
) error {
	return dal4spaceus.RunModuleSpaceWorkerWithUserCtx(ctx, userCtx, request.SpaceID, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo), worker)
}

func RunContactusSpaceWorkerTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCtx facade.UserContext,
	spaceID string,
	worker ContactusModuleWorker,
) error {
	return dal4spaceus.RunModuleSpaceWorkerTx(ctx, tx, userCtx, spaceID, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo), worker)
}

func RunContactusSpaceWorkerNoUpdate(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCtx facade.UserContext,
	spaceID string,
	worker ContactusModuleWorker,
) error {
	return dal4spaceus.RunModuleSpaceWorkerNoUpdates(ctx, tx, userCtx, spaceID, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo), worker)
}
