package dal4scrumus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
)

const ScrumusModuleID = "scrumus"

type ScrumSpaceEntry = record.DataWithID[string, *dbo4scrumus.ScrumSpaceDto]

func NewScrumSpaceKey(id string) *dal.Key {
	key := dal4spaceus.NewSpaceKey(id)
	return dal.NewKeyWithParentAndID(key, dal4spaceus.SpaceModulesCollection, ScrumusModuleID)
}

func NewScrumSpaceEntry(id string) ScrumSpaceEntry {
	key := NewScrumSpaceKey(id)
	return record.NewDataWithID(id, key, new(dbo4scrumus.ScrumSpaceDto))
}

func GetScrumSpace(ctx context.Context, tx dal.ReadTransaction, id string) (ScrumSpaceEntry, error) {
	v := NewScrumSpaceEntry(id)
	return v, tx.Get(ctx, v.Record)
}
