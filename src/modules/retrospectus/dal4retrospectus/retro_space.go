package dal4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
)

const RetrospectusModuleID = "retrospectus"

type RetroSpaceEntry = record.DataWithID[string, *dbo4retrospectus.RetroSpaceDbo]

func NewRetroSpaceKey(id string) *dal.Key {
	key := dbo4spaceus.NewSpaceKey(id)
	return dal.NewKeyWithParentAndID(key, dbo4spaceus.SpaceModulesCollection, RetrospectusModuleID)
}

func NewRetroSpaceEntry(id string) RetroSpaceEntry {
	key := NewRetroSpaceKey(id)
	return record.NewDataWithID(id, key, new(dbo4retrospectus.RetroSpaceDbo))
}

func GetRetroSpaceEntry(ctx context.Context, tx dal.ReadTransaction, id string) (RetroSpaceEntry, error) {
	v := NewRetroSpaceEntry(id)
	return v, tx.Get(ctx, v.Record)
}
