package dal4teamus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
)

type TeamEntry = record.DataWithID[string, *dbo4teamus.TeamDbo]

func NewTeamEntry(id string) (team TeamEntry) {
	teamDto := new(dbo4teamus.TeamDbo)
	return NewTeamEntryWithDto(id, teamDto)
}

func NewTeamEntryWithDto(id string, dto *dbo4teamus.TeamDbo) (team TeamEntry) {
	if dto == nil {
		panic("required parameter dto is nil")
	}
	team = record.NewDataWithID(id, NewTeamKey(id), dto)
	team.ID = id
	return
}
