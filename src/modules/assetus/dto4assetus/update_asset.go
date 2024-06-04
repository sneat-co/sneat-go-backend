package dto4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/extras4assetus"
)

type UpdateAssetRequest struct {
	AssetRequest
	Extra extras4assetus.AssetExtra
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
