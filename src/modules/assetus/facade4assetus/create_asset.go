package facade4assetus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/random"
)

// CreateAssetResponse DTO
type CreateAssetResponse struct {
	ID   string                   `json:"id"`
	Data dbo4assetus.AssetBaseDbo `json:"data"`
}

// CreateAsset creates an asset
func CreateAsset(ctx context.Context, user facade.User, request dto4assetus.CreateAssetRequest) (response CreateAssetResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4teamus.CreateTeamItem(ctx, user,
		request.TeamRequest, const4assetus.ModuleID, new(dbo4assetus.AssetusTeamDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*dbo4assetus.AssetusTeamDbo]) (err error) {
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
	params *dal4teamus.ModuleTeamWorkerParams[*dbo4assetus.AssetusTeamDbo],
) (
	response CreateAssetResponse, err error,
) {
	asset := dal4assetus.NewAssetEntry(request.TeamRequest.TeamID, random.ID(8)) // TODO: use DALgo random ID generator
	asset.Data.AssetBaseDbo = request.Asset
	asset.Data.UserIDs = []string{params.UserID}
	asset.Data.TeamIDs = []string{request.TeamRequest.TeamID}
	//asset.Data.ContactIDs = []string{"*"}
	asset.Data.WithModified = dbmodels.NewWithModified(params.Started, params.UserID)

	asset.Record.SetError(nil) // Mark record as not having an error

	if err = asset.Data.Validate(); err != nil {
		return response, fmt.Errorf("assert record data is not valid before insert: %w", err)
	}

	if err = tx.Set(ctx, asset.Record); err != nil { // TODO: change to .Insert() with random ID generator
		return response, fmt.Errorf("failed to insert asset record: %w", err)
	}

	response.ID = asset.ID

	var assetBrief briefs4assetus.AssetBrief
	if assetBrief, err = asset.Data.GetAssetBrief(); err != nil {
		return
	}

	var assetusTeamModuleUpdates []dal.Update
	if assetusTeamModuleUpdates, err = params.TeamModuleEntry.Data.AddAssetBrief(asset.ID, assetBrief); err != nil {
		return
	}

	if err = params.TeamModuleEntry.Data.Validate(); err != nil {
		return response, fmt.Errorf("assetus team module record is not valid before saving to db: %w", err)
	}

	if params.TeamModuleEntry.Record.Exists() {
		if err = tx.Update(ctx, params.TeamModuleEntry.Record.Key(), assetusTeamModuleUpdates); err != nil {
			return
		}
	} else {
		if err = tx.Insert(ctx, params.TeamModuleEntry.Record); err != nil {
			return
		}
	}

	response.Data = asset.Data.AssetBaseDbo
	return response, err
}
