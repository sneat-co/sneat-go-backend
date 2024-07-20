package dal4spaceus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
)

type SpaceEntry = record.DataWithID[string, *dbo4spaceus.SpaceDbo]

func NewSpaceEntry(id string) (team SpaceEntry) {
	teamDto := new(dbo4spaceus.SpaceDbo)
	return NewSpaceEntryWithDto(id, teamDto)
}

func NewSpaceEntryWithDto(id string, dto *dbo4spaceus.SpaceDbo) (space SpaceEntry) {
	if dto == nil {
		panic("required parameter dto is nil")
	}
	space = record.NewDataWithID(id, NewSpaceKey(id), dto)
	space.ID = id
	return
}
