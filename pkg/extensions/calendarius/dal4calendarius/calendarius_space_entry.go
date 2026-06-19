package dal4calendarius

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/const4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

type CalendariusSpaceEntry = record.DataWithID[coretypes.ExtID, *dbo4calendarius.CalendariusSpaceDbo]

func NewCalendariusSpaceKey(spaceID coretypes.SpaceID) *dal.Key {
	return dbo4spaceus.NewSpaceModuleKey(spaceID, const4calendarius.ExtensionID)
}

func NewCalendariusSpaceEntry(spaceID coretypes.SpaceID) CalendariusSpaceEntry {
	key := NewCalendariusSpaceKey(spaceID)
	return record.NewDataWithID(const4calendarius.ExtensionID, key, new(dbo4calendarius.CalendariusSpaceDbo))
}

func GetCalendariusSpace(ctx context.Context, tx dal.ReadwriteTransaction, spaceID coretypes.SpaceID) (CalendariusSpaceEntry, error) {
	v := NewCalendariusSpaceEntry(spaceID)
	return v, tx.Get(ctx, v.Record)
}
