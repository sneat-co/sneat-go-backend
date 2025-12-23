package facade4assetus

import (
	"context"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// CreateAssetResponse DTO
type CreateAssetResponse struct {
	ID   string                   `json:"id"`
	Data dbo4assetus.AssetBaseDbo `json:"data"`
}

// CreateAsset creates an asset
func CreateAsset(ctx facade.ContextWithUser, request dto4assetus.CreateAssetRequest) (response CreateAssetResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4spaceus.CreateSpaceItem(ctx,
		request.SpaceRequest, const4assetus.ExtensionID, new(dbo4assetus.AssetusSpaceDbo),
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4assetus.AssetusSpaceDbo]) (err error) {
			if err = params.GetRecords(ctx, tx); err != nil {
				return err
			}
			response, err = createAssetTx(ctx, tx, request, params)
			return err
		},
	)
	return
}

func createAssetTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4assetus.CreateAssetRequest,
	params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4assetus.AssetusSpaceDbo],
) (
	response CreateAssetResponse, err error,
) {
	asset := dal4assetus.NewAssetEntryWithoutID(request.SpaceID) // TODO: use DALgo random ContactID generator
	asset.Data.AssetBaseDbo = request.Asset
	asset.Data.UserIDs = []string{params.UserID()}
	asset.Data.SpaceIDs = []coretypes.SpaceID{
		request.SpaceID,
	}
	//asset.Data.ContactIDs = []string{"*"}
	asset.Data.WithModified = dbmodels.NewWithModified(params.Started, params.UserID())

	asset.Record.SetError(nil) // Mark a record as not having an error so we can access record.Data()

	if err = asset.Data.Validate(); err != nil {
		return response, fmt.Errorf("assert record data is not valid before insert: %w", err)
	}

	if err = tx.Insert(ctx, asset.Record, dal.WithRandomStringKey(5, 3)); err != nil {
		return response, fmt.Errorf("failed to insert asset record: %w", err)
	}

	response.ID = asset.ID

	var assetBrief briefs4assetus.AssetBrief
	if assetBrief, err = asset.Data.GetAssetBrief(); err != nil {
		return
	}

	var assetusSpaceModuleUpdates []update.Update
	if assetusSpaceModuleUpdates, err = params.SpaceModuleEntry.Data.AddAssetBrief(asset.ID, assetBrief); err != nil {
		return
	}

	if err = params.SpaceModuleEntry.Data.Validate(); err != nil {
		return response, fmt.Errorf("assetus team module record is not valid before saving to db: %w", err)
	}

	if params.SpaceModuleEntry.Record.Exists() {
		if err = tx.Update(ctx, params.SpaceModuleEntry.Record.Key(), assetusSpaceModuleUpdates); err != nil {
			return
		}
	} else {
		if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
			return
		}
	}

	response.Data = asset.Data.AssetBaseDbo
	return response, err
}
