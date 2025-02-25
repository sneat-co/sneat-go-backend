package dal4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	dbo4spaceus2 "github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

const RetrospectusModuleID coretypes.ModuleID = "retrospectus"

type RetroSpaceEntry = record.DataWithID[coretypes.ModuleID, *dbo4retrospectus.RetroSpaceDbo]

func NewRetroSpaceKey(id coretypes.SpaceID) *dal.Key {
	key := dbo4spaceus2.NewSpaceKey(id)
	return dal.NewKeyWithParentAndID(key, dbo4spaceus2.SpaceModulesCollection, RetrospectusModuleID)
}

func NewRetroSpaceEntry(spaceID coretypes.SpaceID) RetroSpaceEntry {
	key := NewRetroSpaceKey(spaceID)
	return record.NewDataWithID(RetrospectusModuleID, key, new(dbo4retrospectus.RetroSpaceDbo))
}

func GetRetroSpaceEntry(ctx context.Context, tx dal.ReadTransaction, spaceID coretypes.SpaceID) (RetroSpaceEntry, error) {
	v := NewRetroSpaceEntry(spaceID)
	return v, tx.Get(ctx, v.Record)
}
