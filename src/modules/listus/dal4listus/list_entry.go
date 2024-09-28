package dal4listus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
)

type ListEntry = record.DataWithID[string, *dbo4listus.ListDbo]

// NewListKey creates a new list key
func NewListKey(spaceID string, listKey dbo4listus.ListKey) *dal.Key {
	spaceModuleKey := dbo4spaceus.NewSpaceModuleKey(spaceID, const4listus.ModuleID)
	return dal.NewKeyWithParentAndID(spaceModuleKey, dbo4listus.ListsCollection, string(listKey))
}

func NewListEntry(spaceID string, listKey dbo4listus.ListKey) (list ListEntry) {
	key := NewListKey(spaceID, listKey)
	list.ID = key.ID.(string)
	list.FullID = fmt.Sprintf("%s/%s", spaceID, listKey) // TODO: Do we need this?
	list.Key = key
	list.Data = new(dbo4listus.ListDbo)
	list.Record = dal.NewRecordWithData(key, list.Data)
	return
}
