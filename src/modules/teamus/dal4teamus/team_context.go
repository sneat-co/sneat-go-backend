package dal4teamus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
)

type TeamContext = record.DataWithID[string, *models4teamus.TeamDto]

func NewTeamContext(id string) (team TeamContext) {
	teamDto := new(models4teamus.TeamDto)
	return NewTeamContextWithDto(id, teamDto)
}

func NewTeamContextWithDto(id string, dto *models4teamus.TeamDto) (team TeamContext) {
	if dto == nil {
		panic("required parameter dto is nil")
	}
	team = record.NewDataWithID(id, NewTeamKey(id), dto)
	team.ID = id
	return
}
