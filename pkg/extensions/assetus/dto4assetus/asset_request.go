package dto4assetus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

type AssetRequest struct {
	AssetID       string
	AssetCategory string
	dto4spaceus.SpaceRequest
}

func (v AssetRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if err := validate.RecordID(v.AssetID); err != nil {
		return validation.NewErrBadRequestFieldValue("assetID", err.Error())
	}
	return nil
}
