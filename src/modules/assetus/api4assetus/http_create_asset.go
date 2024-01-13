package api4assetus

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/facade4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

// httpPostCreateAsset creates an asset
func httpPostCreateAsset(w http.ResponseWriter, r *http.Request) {
	var request facade4assetus.CreateAssetRequest
	assetCategory := r.URL.Query().Get("assetCategory")
	switch assetCategory {
	case const4assetus.AssetCategoryVehicle:
		asset := models4assetus.NewVehicleAssetDbData()
		asset.Title = asset.GenerateTitle()
		request.Asset = &asset.VehicleAssetMainData
		request.DbData = asset
	//case const4assetus.AssetCategoryRealEstate:
	//	asset := models4assetus.NewDwellingAssetDbData()
	//	request.Asset = &asset.AssetDtoDwelling
	//	request.DbData = asset
	case const4assetus.AssetCategoryDocument:
		asset := models4assetus.NewDocumentDbData()
		request.Asset = asset.DocumentMainData
		request.DbData = asset
	case "":
		apicore.ReturnError(r.Context(), w, r, errors.New("GET parameter 'assetCategory' is required"))
		return
	default:
		apicore.ReturnError(r.Context(), w, r, fmt.Errorf("unsupported asset category: %s", assetCategory))
		return
	}
	handler := func(ctx context.Context, userCtx facade.User) (interface{}, error) {
		asset, err := facade4assetus.CreateAsset(ctx, userCtx, request)
		if err != nil {
			return nil, fmt.Errorf("failed to create asset: %w", err)
		}
		if asset.ID == "" {
			return nil, errors.New("asset created by facade does not have an ContactID")
		}
		if asset.Data == nil {
			return nil, errors.New("asset created by facade does not have a DTO")
		}
		if err = asset.Data.Validate(); err != nil {
			err = fmt.Errorf("asset created by facade is not valid: %w", err)
			return asset, err
		}
		return asset, nil
	}
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, handler, http.StatusCreated, verify.DefaultJsonWithAuthRequired)
}
