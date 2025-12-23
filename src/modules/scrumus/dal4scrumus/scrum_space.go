package dal4scrumus

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	dbo4spaceus2 "github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

const ScrumusModuleID coretypes.ExtID = "scrumus"

type ScrumSpaceEntry = record.DataWithID[coretypes.ExtID, *dbo4scrumus.ScrumSpaceDto]

func NewScrumSpaceKey(id coretypes.SpaceID) *dal.Key {
	key := dbo4spaceus2.NewSpaceKey(id)
	return dal.NewKeyWithParentAndID(key, dbo4spaceus2.SpaceModulesCollection, ScrumusModuleID)
}

func NewScrumSpaceEntry(id coretypes.SpaceID) ScrumSpaceEntry {
	key := NewScrumSpaceKey(id)
	return record.NewDataWithID(ScrumusModuleID, key, new(dbo4scrumus.ScrumSpaceDto))
}

func GetScrumSpace(ctx context.Context, tx dal.ReadTransaction, spaceID coretypes.SpaceID) (ScrumSpaceEntry, error) {
	v := NewScrumSpaceEntry(spaceID)
	return v, tx.Get(ctx, v.Record)
}
