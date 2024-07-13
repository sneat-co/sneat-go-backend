package dal4assetus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type AssetusSpaceWorkerParams = dal4teamus.ModuleSpaceWorkerParams[*dbo4assetus.AssetusSpaceDbo]

func NewAssetusSpaceWorkerParams(userID, spaceID string) *AssetusSpaceWorkerParams {
	spaceWorkerParams := dal4teamus.NewSpaceWorkerParams(userID, spaceID)
	return dal4teamus.NewSpaceModuleWorkerParams(const4assetus.ModuleID, spaceWorkerParams, new(dbo4assetus.AssetusSpaceDbo))
}

func RunReadonlyAssetusSpaceWorker(
	ctx context.Context,
	user facade.User,
	request dto4teamus.SpaceRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, params *AssetusSpaceWorkerParams) (err error),
) error {
	return dal4teamus.RunReadonlyModuleSpaceWorker(ctx, user, request, const4assetus.ModuleID, new(dbo4assetus.AssetusSpaceDbo), worker)
}

type AssetusModuleWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *AssetusSpaceWorkerParams) (err error)

func RunAssetusSpaceWorker(
	ctx context.Context,
	user facade.User,
	request dto4teamus.SpaceRequest,
	worker AssetusModuleWorker,
) error {
	return dal4teamus.RunModuleSpaceWorker(ctx, user, request, const4contactus.ModuleID, new(dbo4assetus.AssetusSpaceDbo), worker)
}

func RunAssetusSpaceWorkerTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	user facade.User,
	request dto4teamus.SpaceRequest,
	worker AssetusModuleWorker,
) error {
	return dal4teamus.RunModuleSpaceWorkerTx(ctx, tx, user, request, const4contactus.ModuleID, new(dbo4assetus.AssetusSpaceDbo), worker)
}
