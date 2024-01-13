package const4assetus

import (
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
)

type AssetStatus = string

const (
	AssetStatusActive   AssetStatus = dbmodels.StatusActive
	AssetStatusArchived AssetStatus = dbmodels.StatusArchived
	AssetStatusDraft    AssetStatus = dbmodels.StatusDraft
)

var AssetStatuses = []AssetStatus{
	AssetStatusActive,
	AssetStatusArchived,
	AssetStatusDraft,
}

func IsValidAssetStatus(status AssetStatus) bool {
	return slice.Contains(AssetStatuses, status)
}
