package dal4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

type CalendariumSpaceEntry = record.DataWithID[string, *dbo4calendarium.CalendariumSpaceDbo]

func NewCalendariumSpaceKey(spaceID string) *dal.Key {
	return dal4teamus.NewSpaceModuleKey(spaceID, const4calendarium.ModuleID)
}

func NewCalendariumSpaceEntry(spaceID string) CalendariumSpaceEntry {
	key := NewCalendariumSpaceKey(spaceID)
	return record.NewDataWithID(spaceID, key, new(dbo4calendarium.CalendariumSpaceDbo))
}

func GetCalendariumSpace(ctx context.Context, tx dal.ReadwriteTransaction, spaceID string) (CalendariumSpaceEntry, error) {
	v := NewCalendariumSpaceEntry(spaceID)
	return v, tx.Get(ctx, v.Record)
}
