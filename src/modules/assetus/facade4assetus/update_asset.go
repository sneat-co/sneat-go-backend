package facade4assetus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateAsset(ctx context.Context, user facade.User, request dto4assetus.UpdateAssetRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return UpdateAssetTx(ctx, tx, user, request)
	})
}

func UpdateAssetTx(ctx context.Context, tx dal.ReadwriteTransaction, user facade.User, request dto4assetus.UpdateAssetRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	switch request.AssetCategory {
	case "vehicle":
		return runAssetWorker(ctx, tx, user, request, new(models4assetus.AssetVehicleExtra))
	case "dwelling":
		return runAssetWorker(ctx, tx, user, request, new(models4assetus.AssetDwellingExtra))
	default:
		return runAssetWorker(ctx, tx, user, request, models4assetus.NewAssetNoExtra())
	}
}

type AssetWorkerParams struct {
	*dal4teamus.ModuleTeamWorkerParams[*models4assetus.AssetusTeamDto]
	Asset        record.DataWithID[string, *models4assetus.AssetDbo]
	AssetUpdates []dal.Update
}

func runAssetWorker(ctx context.Context, tx dal.ReadwriteTransaction, user facade.User, request dto4assetus.UpdateAssetRequest, assetExtra models4assetus.AssetExtra) (err error) {
	// TODO: Replace with future RunTeamModuleItemWorkerTx
	return dal4teamus.RunModuleTeamWorkerTx[*models4assetus.AssetusTeamDto](ctx, tx, user, request.TeamRequest, const4assetus.ModuleID, new(models4assetus.AssetusTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *dal4teamus.ModuleTeamWorkerParams[*models4assetus.AssetusTeamDto]) (err error) {
			params := AssetWorkerParams{
				Asset:                  NewAsset("", assetExtra),
				ModuleTeamWorkerParams: teamWorkerParams,
			}
			if err := tx.Get(ctx, params.Asset.Record); err != nil {
				return err
			}
			if err = updateAssetTxWorker(ctx, tx, request, &params); err != nil {
				return err
			}
			if len(params.AssetUpdates) > 0 {
				if err = params.Asset.Data.Validate(); err != nil {
					return fmt.Errorf("asset data is not valid before updating asset record: %w", err)
				}
				if err = tx.Update(ctx, params.Asset.Key, params.AssetUpdates); err != nil {
					return err
				}
			}
			return err
		},
	)
}

func updateAssetTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, _ dto4assetus.UpdateAssetRequest, params *AssetWorkerParams) (err error) {
	if err = tx.Get(ctx, params.Asset.Record); err != nil {
		return fmt.Errorf("failed to get asset record: %w", err)
	}

	if err := params.Asset.Data.Validate(); err != nil {
		return fmt.Errorf("asset DBO is not valid after loading from DB: %w", err)
	}
	return err
}
