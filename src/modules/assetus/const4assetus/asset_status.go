package const4assetus

import (
	"slices"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
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
