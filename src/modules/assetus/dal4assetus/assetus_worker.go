package dal4assetus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	dal4spaceus2 "github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

type AssetusSpaceWorkerParams = dal4spaceus2.ModuleSpaceWorkerParams[*dbo4assetus.AssetusSpaceDbo]

func NewAssetusSpaceWorkerParams(userCtx facade.UserContext, spaceID coretypes.SpaceID) *AssetusSpaceWorkerParams {
	spaceWorkerParams := dal4spaceus2.NewSpaceWorkerParams(userCtx, spaceID)
	return dal4spaceus2.NewSpaceModuleWorkerParams(const4assetus.ExtensionID, spaceWorkerParams, new(dbo4assetus.AssetusSpaceDbo))
}

func RunReadonlyAssetusSpaceWorker(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, params *AssetusSpaceWorkerParams) (err error),
) error {
	return dal4spaceus2.RunReadonlyModuleSpaceWorker(ctx, userCtx, request, const4assetus.ExtensionID, new(dbo4assetus.AssetusSpaceDbo), worker)
}

type AssetusModuleWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *AssetusSpaceWorkerParams) (err error)

func RunAssetusSpaceWorker(
	ctx facade.ContextWithUser,
	request dto4spaceus.SpaceRequest,
	worker AssetusModuleWorker,
) error {
	return dal4spaceus2.RunModuleSpaceWorkerWithUserCtx(ctx, request.SpaceID, const4contactus.ExtensionID, new(dbo4assetus.AssetusSpaceDbo), worker)
}

func RunAssetusSpaceWorkerTx(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	spaceID coretypes.SpaceID,
	worker AssetusModuleWorker,
) error {
	return dal4spaceus2.RunModuleSpaceWorkerTx(ctx, tx, spaceID, const4contactus.ExtensionID, new(dbo4assetus.AssetusSpaceDbo), worker)
}
