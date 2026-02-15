package api4assetus

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

// httpPostCreateAsset creates an asset
func httpPostCreateAsset(w http.ResponseWriter, r *http.Request) {
	var (
		request dto4assetus.CreateAssetRequest
		err     error
	)

	// Create asset base DBO with a specific extra data based on the asset category
	assetCategory := r.URL.Query().Get("assetCategory")
	if request.Asset, err = createAssetBaseDbo(assetCategory); err != nil {
		apicore.ReturnError(r.Context(), w, r, err)
	}

	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			asset, err := facade4assetus.CreateAsset(ctx, request)
			if err != nil {
				return nil, fmt.Errorf("failed to create asset: %w", err)
			}
			if asset.ID == "" {
				return nil, errors.New("created asset does not have an ID")
			}
			if err = asset.Data.Validate(); err != nil {
				err = fmt.Errorf("created asset is not valid: %w", err)
				return asset, err
			}
			return asset, nil
		},
	)
}
