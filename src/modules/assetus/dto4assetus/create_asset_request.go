package dto4assetus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
)

type CreateAssetRequest struct {
	dto4spaceus.SpaceRequest
	Asset dbo4assetus.AssetBaseDbo `json:"asset"`
}
