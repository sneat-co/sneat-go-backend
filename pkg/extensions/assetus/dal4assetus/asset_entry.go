package dal4assetus

import (
	"fmt"
	"reflect"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/dbo4assetus"
	core "github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func NewAssetEntryWithoutID(spaceID coretypes.SpaceID) (asset dbo4assetus.AssetEntry) {
	key := NewAssetKeyWithoutID(spaceID)
	return newAssetEntryWithKey(key)
}

func NewAssetEntry(spaceID coretypes.SpaceID, assetID string) (asset dbo4assetus.AssetEntry) {
	asset.ID = assetID
	asset.Key = NewAssetKey(spaceID, assetID)
	asset.FullID = string(spaceID) + ":" + assetID
	entry := newAssetEntryWithKey(asset.Key)
	entry.ID = asset.ID
	entry.FullID = asset.FullID
	return entry
}

func newAssetEntryWithKey(key *dal.Key) (asset dbo4assetus.AssetEntry) {
	asset.Key = key
	asset.Data = new(dbo4assetus.AssetDbo)
	asset.Record = dal.NewRecordWithData(key, asset.Data)
	return
}

func NewAssetKeyWithoutID(spaceID coretypes.SpaceID) *dal.Key {
	spaceModuleKey := dbo4spaceus.NewSpaceModuleKey(spaceID, const4assetus.ExtensionID)
	return dal.NewIncompleteKey(dbo4assetus.SpaceAssetsCollection, reflect.String, spaceModuleKey)
}

func NewAssetKey(spaceID coretypes.SpaceID, assetID string) *dal.Key {
	if !core.IsAlphanumericOrUnderscore(assetID) {
		panic(fmt.Errorf("assetID should be alphanumeric, got: [%s]", assetID))
	}
	spaceModuleKey := dbo4spaceus.NewSpaceModuleKey(spaceID, const4assetus.ExtensionID)
	return dal.NewKeyWithParentAndID(spaceModuleKey, dbo4assetus.SpaceAssetsCollection, assetID)
}
