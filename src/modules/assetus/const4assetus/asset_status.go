package const4assetus

import (
	"slices"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

type AssetStatus = string

const (
	AssetStatusActive   = dbmodels.StatusActive
	AssetStatusArchived = dbmodels.StatusArchived
	AssetStatusDraft    = dbmodels.StatusDraft
)

var AssetStatuses = []AssetStatus{
	AssetStatusActive,
	AssetStatusArchived,
	AssetStatusDraft,
}

func IsValidAssetStatus(status AssetStatus) bool {
	return slices.Contains(AssetStatuses, status)
}
