package facade4assetus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/random"
)

// AssetSummary DTO
type AssetSummary struct {
	RegNumber      string `json:"number,omitempty"`
	DateOfBuild    string `json:"dateOfBuild,omitempty"`
	DateOfPurchase string `json:"dateOfPurchase,omitempty"`
}

type CreateAssetData struct {
	models4assetus.AssetBaseDbo
}

type AssetDto struct {
}

// CreateAssetRequest is a DTO for creating an asset
type CreateAssetRequest struct {
	dto4teamus.TeamRequest
	Asset   models4assetus.AssetCreationData `json:"asset"`
	Related models4linkage.WithRelated       `json:"related"`
}

// Validate returns error if not valid
func (v CreateAssetRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if err := v.Asset.Validate(); err != nil {
		return err
	}
	return nil
}

// CreateAssetResponse DTO
type CreateAssetResponse struct {
	ID   string                      `json:"id"`
	Data models4assetus.AssetBaseDbo `json:"data"`
}

// CreateAsset creates an asset
func CreateAsset(ctx context.Context, user facade.User, request CreateAssetRequest) (response CreateAssetResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4teamus.CreateTeamItem(ctx, user, "assets", request.TeamRequest, const4assetus.ModuleID, new(models4assetus.AssetusTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4assetus.AssetusTeamDto]) (err error) {
			response, err = createAssetTx(ctx, tx, request, params)
			return err
		},
	)
	return
}

func createAssetTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request CreateAssetRequest,
	params *dal4teamus.ModuleTeamWorkerParams[*models4assetus.AssetusTeamDto],
) (
	response CreateAssetResponse, err error,
) {
	response.ID = random.ID(7) // TODO: consider using incomplete key with options?
	asset := NewAsset("", nil)
	asset.Data.UserIDs = []string{params.UserID}
	asset.Data.TeamIDs = []string{request.TeamRequest.TeamID}
	asset.Data.WithModified = dbmodels.NewWithModified(params.Started, params.UserID)
	if err = tx.Insert(ctx, asset.Record); err != nil {
		return response, fmt.Errorf("failed to insert response record")
	}

	assetusTeamModuleUpdates := params.TeamModuleEntry.Data.WithAssets.AddAsset(response.ID, &asset.Data.AssetBrief)

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

	return response, err
}
