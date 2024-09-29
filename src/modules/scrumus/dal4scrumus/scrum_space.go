package dal4scrumus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	dbo4spaceus2 "github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
)

const ScrumusModuleID = "scrumus"

type ScrumSpaceEntry = record.DataWithID[string, *dbo4scrumus.ScrumSpaceDto]

func NewScrumSpaceKey(id string) *dal.Key {
	key := dbo4spaceus2.NewSpaceKey(id)
	return dal.NewKeyWithParentAndID(key, dbo4spaceus2.SpaceModulesCollection, ScrumusModuleID)
}

func NewScrumSpaceEntry(id string) ScrumSpaceEntry {
	key := NewScrumSpaceKey(id)
	return record.NewDataWithID(id, key, new(dbo4scrumus.ScrumSpaceDto))
}

func GetScrumSpace(ctx context.Context, tx dal.ReadTransaction, id string) (ScrumSpaceEntry, error) {
	v := NewScrumSpaceEntry(id)
	return v, tx.Get(ctx, v.Record)
}
