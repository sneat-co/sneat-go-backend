package facade4assetus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type Asset = record.DataWithID[string, models4assetus.AssetDbData]

func NewAsset(id string, dto models4assetus.AssetDbData) Asset {
	key := dal.NewKeyWithID(dal4assetus.AssetsCollection, id)
	return record.NewDataWithID(id, key, dto)
}

// GetAssetByID returns asset by ID
func GetAssetByID(ctx context.Context, id string, dto models4assetus.AssetDbData) (asset Asset, err error) {
	asset = NewAsset(id, dto)
	db := facade.GetDatabase(ctx)
	return asset, db.Get(ctx, asset.Record)
}
