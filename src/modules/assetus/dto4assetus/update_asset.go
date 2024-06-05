package dto4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodels/extra"
)

type UpdateAssetRequest struct {
	AssetRequest
	Extra extra.Data
}

func (v UpdateAssetRequest) Validate() error {
	if err := v.AssetRequest.Validate(); err != nil {
		return err
	}
	if err := v.Extra.Validate(); err != nil {
		return err
	}
	return nil
}
