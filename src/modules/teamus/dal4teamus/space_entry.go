package dal4teamus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
)

type SpaceEntry = record.DataWithID[string, *dbo4teamus.SpaceDbo]

func NewSpaceEntry(id string) (team SpaceEntry) {
	teamDto := new(dbo4teamus.SpaceDbo)
	return NewSpaceEntryWithDto(id, teamDto)
}

func NewSpaceEntryWithDto(id string, dto *dbo4teamus.SpaceDbo) (space SpaceEntry) {
	if dto == nil {
		panic("required parameter dto is nil")
	}
	space = record.NewDataWithID(id, NewSpaceKey(id), dto)
	space.ID = id
	return
}
