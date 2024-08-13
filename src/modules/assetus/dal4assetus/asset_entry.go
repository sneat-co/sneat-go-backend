package dal4assetus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	core "github.com/sneat-co/sneat-go-core"
)

func NewAssetEntry(teamID, assetID string) (asset dbo4assetus.AssetEntry) {
	key := NewAssetKey(teamID, assetID)
	asset.ID = assetID
	asset.FullID = teamID + ":" + assetID
	asset.Key = key
	asset.Data = new(dbo4assetus.AssetDbo)
	asset.Record = dal.NewRecordWithData(key, asset.Data)
	return
}

func NewAssetKey(teamID, assetID string) *dal.Key {
	if !core.IsAlphanumericOrUnderscore(assetID) {
		panic(fmt.Errorf("assetID should be alphanumeric, got: [%s]", assetID))
	}
	teamModuleKey := dbo4spaceus.NewSpaceModuleKey(teamID, const4assetus.ModuleID)
	return dal.NewKeyWithParentAndID(teamModuleKey, dbo4assetus.SpaceAssetsCollection, assetID)
}
