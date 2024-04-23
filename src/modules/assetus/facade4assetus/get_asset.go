package facade4assetus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
)

type Asset = record.DataWithID[string, *models4assetus.AssetDbo]

func NewAsset(id string, extra models4assetus.AssetExtra) Asset {
	key := dal.NewKeyWithID(dal4assetus.AssetsCollection, id)
	dbo := new(models4assetus.AssetDbo)
	dbo.Extra = extra
	return record.NewDataWithID(id, key, dbo)
}
