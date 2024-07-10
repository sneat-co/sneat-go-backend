package dal4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

const RetrospectusModuleID = "retrospectus"

type RetroSpaceEntry = record.DataWithID[string, *dbo4retrospectus.RetroSpaceDbo]

func NewRetroSpaceKey(id string) *dal.Key {
	key := dal4teamus.NewSpaceKey(id)
	return dal.NewKeyWithParentAndID(key, dal4teamus.SpaceModulesCollection, RetrospectusModuleID)
}

func NewRetroSpaceEntry(id string) RetroSpaceEntry {
	key := NewRetroSpaceKey(id)
	return record.NewDataWithID(id, key, new(dbo4retrospectus.RetroSpaceDbo))
}

func GetRetroSpaceEntry(ctx context.Context, tx dal.ReadTransaction, id string) (RetroSpaceEntry, error) {
	v := NewRetroSpaceEntry(id)
	return v, tx.Get(ctx, v.Record)
}
