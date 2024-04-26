package dto4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
)

type CreateAssetRequest struct {
	dto4teamus.TeamRequest
	Asset models4assetus.AssetBaseDbo `json:"asset"`
}
