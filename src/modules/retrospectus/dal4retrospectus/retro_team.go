package dal4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

const RetrospectusModuleID = "retrospectus"

type RetroTeam = record.DataWithID[string, *dbo4retrospectus.RetroTeamDto]

func NewRetroTeamKey(id string) *dal.Key {
	teamKey := dal4teamus.NewTeamKey(id)
	return dal.NewKeyWithParentAndID(teamKey, dal4teamus.TeamModulesCollection, RetrospectusModuleID)
}

func NewRetroTeam(id string) RetroTeam {
	key := NewRetroTeamKey(id)
	return record.NewDataWithID(id, key, new(dbo4retrospectus.RetroTeamDto))
}

func GetRetroTeam(ctx context.Context, tx dal.ReadTransaction, id string) (RetroTeam, error) {
	retro := NewRetroTeam(id)
	return retro, tx.Get(ctx, retro.Record)
}
