package api4assetus

import (
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"net/http"
)

func createAssetBaseDbo(r *http.Request) (assetDbo models4assetus.AssetBaseDbo, err error) {
	assetCategory := (models4assetus.AssetExtraType)(r.URL.Query().Get("assetCategory"))
	if assetCategory == "" {
		err = errors.New("GET parameter 'assetCategory' is required")
		return
	}
	assetExtra := models4assetus.NewAssetExtra(assetCategory)
	if assetExtra == nil {
		err = fmt.Errorf("unsupported asset category: %s", assetCategory)
		return
	}
	if err = assetDbo.SetExtra(assetExtra); err != nil {
		err = fmt.Errorf("failed to set asset extra data: %w", err)
		return
	}
	return
}
