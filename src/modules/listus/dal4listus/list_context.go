package dal4listus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

type ListContext struct {
	record.WithID[string]
	Dto *dbo4listus.ListDto
}

// NewTeamListKey creates new list key
func NewTeamListKey(teamID, id string) *dal.Key {
	teamModuleKey := dal4teamus.NewTeamModuleKey(teamID, const4listus.ModuleID)
	return dal.NewKeyWithParentAndID(teamModuleKey, dbo4listus.ListsCollection, id)
}

func NewTeamListContext(teamID, listID string) (list ListContext) {
	key := NewTeamListKey(teamID, listID)
	list.ID = listID
	list.FullID = teamID + dbo4listus.ListIDSeparator + listID
	list.Key = key
	list.Dto = new(dbo4listus.ListDto)
	list.Record = dal.NewRecordWithData(key, list.Dto)
	return
}
