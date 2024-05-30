package api4assetus

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/facade4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

// httpPostCreateAsset creates an asset
func httpPostCreateAsset(w http.ResponseWriter, r *http.Request) {
	var request dto4assetus.CreateAssetRequest
	assetCategory := r.URL.Query().Get("assetCategory")
	var assetExtra models4assetus.AssetExtra
	switch assetCategory {
	case const4assetus.AssetCategoryVehicle:
		assetExtra = new(models4assetus.AssetVehicleExtra)
	case const4assetus.AssetCategoryDocument:
		assetExtra = new(models4assetus.AssetDocumentExtra)
	case const4assetus.AssetCategoryDwelling:
		assetExtra = new(models4assetus.AssetDwellingExtra)
	case "":
		apicore.ReturnError(r.Context(), w, r, errors.New("GET parameter 'assetCategory' is required"))
		return
	default:
		apicore.ReturnError(r.Context(), w, r, fmt.Errorf("unsupported asset category: %s", assetCategory))
		return
	}
	if err := request.Asset.SetExtra(assetExtra); err != nil {
		apicore.ReturnError(r.Context(), w, r, fmt.Errorf("failed to set asset extra data: %w", err))
		return
	}
	createAssetHttpHandler := func(ctx context.Context, userCtx facade.User) (interface{}, error) {
		asset, err := facade4assetus.CreateAsset(ctx, userCtx, request)
		if err != nil {
			return nil, fmt.Errorf("failed to create asset: %w", err)
		}
		if asset.ID == "" {
			return nil, errors.New("asset created by facade does not have an ContactID")
		}
		if err = asset.Data.Validate(); err != nil {
			err = fmt.Errorf("asset created by facade is not valid: %w", err)
			return asset, err
		}
		return asset, nil
	}
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, createAssetHttpHandler, http.StatusCreated, verify.DefaultJsonWithAuthRequired)
}
