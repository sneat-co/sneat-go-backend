package dal4listus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
)

type ListEntry = record.DataWithID[string, *dbo4listus.ListDbo]

// NewSpaceListKey creates a new list key
func NewSpaceListKey(teamID, id string) *dal.Key {
	key := dal4spaceus.NewSpaceModuleKey(teamID, const4listus.ModuleID)
	return dal.NewKeyWithParentAndID(key, dbo4listus.ListsCollection, id)
}

func NewSpaceListEntry(teamID, listID string) (list ListEntry) {
	key := NewSpaceListKey(teamID, listID)
	list.ID = listID
	list.FullID = teamID + dbo4listus.ListIDSeparator + listID
	list.Key = key
	list.Data = new(dbo4listus.ListDbo)
	list.Record = dal.NewRecordWithData(key, list.Data)
	return
}
