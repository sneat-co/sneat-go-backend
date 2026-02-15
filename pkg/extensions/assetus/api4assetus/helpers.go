package api4assetus

import (
	"errors"
	"fmt"

	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/extras4assetus"
)

func createAssetBaseDbo(assetCategory string) (assetDbo dbo4assetus.AssetBaseDbo, err error) {
	if assetCategory == "" {
		err = errors.New("GET parameter 'assetCategory' is required")
		return
	}
	extraType := (extra.Type)(assetCategory)

	assetExtra := extras4assetus.NewAssetExtra(extraType)
	if assetExtra == nil {
		err = fmt.Errorf("unsupported asset extra type: %s", extraType)
		return
	}

	if err = assetDbo.SetExtra(extraType, assetExtra); err != nil {
		err = fmt.Errorf("failed to set asset extra data: %w", err)
		return
	}
	return
}
