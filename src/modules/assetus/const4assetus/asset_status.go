package const4assetus

import (
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"slices"
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
	return slices.Contains(AssetStatuses, status)
}
