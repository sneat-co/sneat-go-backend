package dto4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
)

type CreateAssetRequest struct {
	dto4spaceus.SpaceRequest
	Asset dbo4assetus.AssetBaseDbo `json:"asset"`
}
