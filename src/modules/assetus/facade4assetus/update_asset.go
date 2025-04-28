package facade4assetus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	extra2 "github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateAsset(ctx facade.ContextWithUser, request dto4assetus.UpdateAssetRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	return dal4assetus.RunAssetusSpaceWorker(ctx, request.SpaceRequest, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetusSpaceWorkerParams) (err error) {
		err = UpdateAssetTx(ctx, tx, request)
		if err != nil {
			return err
		}
		return errors.New("UpdateAssetTx need to use params argument")
	})
}

func UpdateAssetTx(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, request dto4assetus.UpdateAssetRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	extraData := extra2.NewExtraData(extra2.Type(request.AssetCategory))
	return runAssetWorker(ctx, tx, request, extraData)
}

type AssetWorkerParams struct {
	*dal4spaceus.ModuleSpaceWorkerParams[*dbo4assetus.AssetusSpaceDbo]
	Asset        record.DataWithID[string, *dbo4assetus.AssetDbo]
	AssetUpdates []update.Update
}

func runAssetWorker(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, request dto4assetus.UpdateAssetRequest, extraData extra2.Data) (err error) {
	// TODO: Replace with future RunTeamModuleItemWorkerTx
	return dal4spaceus.RunModuleSpaceWorkerTx[*dbo4assetus.AssetusSpaceDbo](ctx, tx, request.SpaceID, const4assetus.ModuleID, new(dbo4assetus.AssetusSpaceDbo),
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, teamWorkerParams *dal4spaceus.ModuleSpaceWorkerParams[*dbo4assetus.AssetusSpaceDbo]) (err error) {
			extraType := extra2.Type(request.AssetCategory)
			params := AssetWorkerParams{
				Asset:                   NewAsset("", extraType, extraData),
				ModuleSpaceWorkerParams: teamWorkerParams,
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
