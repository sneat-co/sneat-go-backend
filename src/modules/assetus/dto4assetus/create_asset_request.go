package dto4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
)

type CreateAssetRequest struct {
	dto4teamus.TeamRequest
	Asset dbo4assetus.AssetBaseDbo `json:"asset"`
}
