package dal4scrumus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

const ScrumusModuleID = "scrumus"

type ScrumTeam = record.DataWithID[string, *dbo4scrumus.ScrumTeamDto]

func NewScrumTeamKey(id string) *dal.Key {
	teamKey := dal4teamus.NewTeamKey(id)
	return dal.NewKeyWithParentAndID(teamKey, dal4teamus.TeamModulesCollection, ScrumusModuleID)
}

func NewScrumTeam(id string) ScrumTeam {
	key := NewScrumTeamKey(id)
	return record.NewDataWithID(id, key, new(dbo4scrumus.ScrumTeamDto))
}

func GetScrumTeam(ctx context.Context, tx dal.ReadTransaction, id string) (ScrumTeam, error) {
	retro := NewScrumTeam(id)
	return retro, tx.Get(ctx, retro.Record)
}
