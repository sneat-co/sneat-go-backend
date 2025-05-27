package dal4assetus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	core "github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"reflect"
)

func NewAssetEntryWithoutID(spaceID coretypes.SpaceID) (asset dbo4assetus.AssetEntry) {
	asset.Key = NewAssetKeyWithoutID(spaceID)
	return newAssetEntryWithKey(asset.Key)
}
func NewAssetEntry(spaceID coretypes.SpaceID, assetID string) (asset dbo4assetus.AssetEntry) {
	asset.ID = assetID
	asset.Key = NewAssetKey(spaceID, assetID)
	asset.FullID = string(spaceID) + ":" + assetID
	return newAssetEntryWithKey(asset.Key)
}

func newAssetEntryWithKey(key *dal.Key) (asset dbo4assetus.AssetEntry) {
	asset.Key = key
	asset.Data = new(dbo4assetus.AssetDbo)
	asset.Record = dal.NewRecordWithData(key, asset.Data)
	return
}

func NewAssetKeyWithoutID(spaceID coretypes.SpaceID) *dal.Key {
	spaceModuleKey := dbo4spaceus.NewSpaceModuleKey(spaceID, const4assetus.ModuleID)
	return dal.NewIncompleteKey(dbo4assetus.SpaceAssetsCollection, reflect.String, spaceModuleKey)
}

func NewAssetKey(spaceID coretypes.SpaceID, assetID string) *dal.Key {
	if !core.IsAlphanumericOrUnderscore(assetID) {
		panic(fmt.Errorf("assetID should be alphanumeric, got: [%s]", assetID))
	}
	spaceModuleKey := dbo4spaceus.NewSpaceModuleKey(spaceID, const4assetus.ModuleID)
	return dal.NewKeyWithParentAndID(spaceModuleKey, dbo4assetus.SpaceAssetsCollection, assetID)
}
