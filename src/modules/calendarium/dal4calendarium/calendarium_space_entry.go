package dal4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
)

type CalendariumSpaceEntry = record.DataWithID[string, *dbo4calendarium.CalendariumSpaceDbo]

func NewCalendariumSpaceKey(spaceID string) *dal.Key {
	return dbo4spaceus.NewSpaceModuleKey(spaceID, const4calendarium.ModuleID)
}

func NewCalendariumSpaceEntry(spaceID string) CalendariumSpaceEntry {
	key := NewCalendariumSpaceKey(spaceID)
	return record.NewDataWithID(spaceID, key, new(dbo4calendarium.CalendariumSpaceDbo))
}

func GetCalendariumSpace(ctx context.Context, tx dal.ReadwriteTransaction, spaceID string) (CalendariumSpaceEntry, error) {
	v := NewCalendariumSpaceEntry(spaceID)
	return v, tx.Get(ctx, v.Record)
}
