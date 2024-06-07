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

type AssetusTeamWorkerParams = dal4teamus.ModuleTeamWorkerParams[*dbo4assetus.AssetusTeamDbo]

func NewAssetusTeamWorkerParams(userID, teamID string) *AssetusTeamWorkerParams {
	teamWorkerParams := dal4teamus.NewTeamWorkerParams(userID, teamID)
	return dal4teamus.NewTeamModuleWorkerParams(const4assetus.ModuleID, teamWorkerParams, new(dbo4assetus.AssetusTeamDbo))
}

func RunReadonlyAssetusTeamWorker(
	ctx context.Context,
	user facade.User,
	request dto4teamus.TeamRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, params *AssetusTeamWorkerParams) (err error),
) error {
	return dal4teamus.RunReadonlyModuleTeamWorker(ctx, user, request, const4assetus.ModuleID, new(dbo4assetus.AssetusTeamDbo), worker)
}

type AssetusModuleWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *AssetusTeamWorkerParams) (err error)

func RunAssetusTeamWorker(
	ctx context.Context,
	user facade.User,
	request dto4teamus.TeamRequest,
	worker AssetusModuleWorker,
) error {
	return dal4teamus.RunModuleTeamWorker(ctx, user, request, const4contactus.ModuleID, new(dbo4assetus.AssetusTeamDbo), worker)
}

func RunAssetusTeamWorkerTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	user facade.User,
	request dto4teamus.TeamRequest,
	worker AssetusModuleWorker,
) error {
	return dal4teamus.RunModuleTeamWorkerTx(ctx, tx, user, request, const4contactus.ModuleID, new(dbo4assetus.AssetusTeamDbo), worker)
}
