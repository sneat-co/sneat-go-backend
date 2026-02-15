package facade4assetus

import (
	"fmt"
	"reflect"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/dal4assetus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/dbo4assetus"
)

type Asset = record.DataWithID[string, *dbo4assetus.AssetDbo]

func NewAsset(id string, extraType extra.Type, extraData extra.Data) Asset {
	var key *dal.Key
	if id == "" {
		key = dal.NewIncompleteKey(dal4assetus.AssetsCollection, reflect.String, nil)
	} else {
		key = dal.NewKeyWithID(dal4assetus.AssetsCollection, id)
	}
	dbo := new(dbo4assetus.AssetDbo)
	if err := dbo.SetExtra(extraType, extraData); err != nil {
		panic(fmt.Errorf("failed to set asset extra data: %w", err))
	}
	return record.NewDataWithID(id, key, dbo)
}
