package dal4assetus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type AssetusSpaceWorkerParams = dal4spaceus.ModuleSpaceWorkerParams[*dbo4assetus.AssetusSpaceDbo]

func NewAssetusSpaceWorkerParams(userCtx facade.UserContext, spaceID string) *AssetusSpaceWorkerParams {
	spaceWorkerParams := dal4spaceus.NewSpaceWorkerParams(userCtx, spaceID)
	return dal4spaceus.NewSpaceModuleWorkerParams(const4assetus.ModuleID, spaceWorkerParams, new(dbo4assetus.AssetusSpaceDbo))
}

func RunReadonlyAssetusSpaceWorker(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, params *AssetusSpaceWorkerParams) (err error),
) error {
	return dal4spaceus.RunReadonlyModuleSpaceWorker(ctx, userCtx, request, const4assetus.ModuleID, new(dbo4assetus.AssetusSpaceDbo), worker)
}

type AssetusModuleWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *AssetusSpaceWorkerParams) (err error)

func RunAssetusSpaceWorker(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	worker AssetusModuleWorker,
) error {
	return dal4spaceus.RunModuleSpaceWorkerWithUserCtx(ctx, userCtx, request.SpaceID, const4contactus.ModuleID, new(dbo4assetus.AssetusSpaceDbo), worker)
}

func RunAssetusSpaceWorkerTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCtx facade.UserContext,
	spaceID string,
	worker AssetusModuleWorker,
) error {
	return dal4spaceus.RunModuleSpaceWorkerTx(ctx, tx, userCtx, spaceID, const4contactus.ModuleID, new(dbo4assetus.AssetusSpaceDbo), worker)
}
