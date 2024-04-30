package dal4assetus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	core "github.com/sneat-co/sneat-go-core"
)

func NewAssetEntry(teamID, assetID string) (asset record.DataWithID[string, *models4assetus.AssetDbo]) {
	key := NewAssetKey(teamID, assetID)
	asset.ID = assetID
	asset.FullID = teamID + ":" + assetID
	asset.Key = key
	asset.Data = new(models4assetus.AssetDbo)
	asset.Record = dal.NewRecordWithData(key, asset.Data)
	return
}

func NewAssetKey(teamID, assetID string) *dal.Key {
	if !core.IsAlphanumericOrUnderscore(assetID) {
		panic(fmt.Errorf("assetID should be alphanumeric, got: [%v]", assetID))
	}
	teamModuleKey := dal4teamus.NewTeamModuleKey(teamID, const4assetus.ModuleID)
	return dal.NewKeyWithParentAndID(teamModuleKey, models4assetus.TeamAssetsCollection, assetID)
}
