package dal4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	dbo4spaceus2 "github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
)

const RetrospectusModuleID = "retrospectus"

type RetroSpaceEntry = record.DataWithID[string, *dbo4retrospectus.RetroSpaceDbo]

func NewRetroSpaceKey(id string) *dal.Key {
	key := dbo4spaceus2.NewSpaceKey(id)
	return dal.NewKeyWithParentAndID(key, dbo4spaceus2.SpaceModulesCollection, RetrospectusModuleID)
}

func NewRetroSpaceEntry(id string) RetroSpaceEntry {
	key := NewRetroSpaceKey(id)
	return record.NewDataWithID(id, key, new(dbo4retrospectus.RetroSpaceDbo))
}

func GetRetroSpaceEntry(ctx context.Context, tx dal.ReadTransaction, id string) (RetroSpaceEntry, error) {
	v := NewRetroSpaceEntry(id)
	return v, tx.Get(ctx, v.Record)
}
