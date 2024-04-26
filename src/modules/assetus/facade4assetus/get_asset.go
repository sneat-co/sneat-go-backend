package facade4assetus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"reflect"
)

type Asset = record.DataWithID[string, *models4assetus.AssetDbo]

func NewAsset(id string, extra models4assetus.AssetExtra) Asset {
	var key *dal.Key
	if id == "" {
		key = dal.NewIncompleteKey(dal4assetus.AssetsCollection, reflect.String, nil)
	} else {
		key = dal.NewKeyWithID(dal4assetus.AssetsCollection, id)
	}
	dbo := new(models4assetus.AssetDbo)
	if err := dbo.SetExtra(extra); err != nil {
		panic(fmt.Errorf("failed to set asset extra data: %w", err))
	}
	return record.NewDataWithID(id, key, dbo)
}
