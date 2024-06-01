package api4assetus

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

// httpPostCreateAsset creates an asset
func httpPostCreateAsset(w http.ResponseWriter, r *http.Request) {
	var request dto4assetus.CreateAssetRequest
	assetDbo, err := createAssetBaseDbo(r)
	if err != nil {
		apicore.ReturnError(r.Context(), w, r, err)
	}
	request.Asset = assetDbo
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
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
		})
}
