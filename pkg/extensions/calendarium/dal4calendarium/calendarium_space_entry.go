package dal4calendarium

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

type CalendariumSpaceEntry = record.DataWithID[coretypes.ExtID, *dbo4calendarium.CalendariumSpaceDbo]

func NewCalendariumSpaceKey(spaceID coretypes.SpaceID) *dal.Key {
	return dbo4spaceus.NewSpaceModuleKey(spaceID, const4calendarium.ExtensionID)
}

func NewCalendariumSpaceEntry(spaceID coretypes.SpaceID) CalendariumSpaceEntry {
	key := NewCalendariumSpaceKey(spaceID)
	return record.NewDataWithID(const4calendarium.ExtensionID, key, new(dbo4calendarium.CalendariumSpaceDbo))
}

func GetCalendariumSpace(ctx context.Context, tx dal.ReadwriteTransaction, spaceID coretypes.SpaceID) (CalendariumSpaceEntry, error) {
	v := NewCalendariumSpaceEntry(spaceID)
	return v, tx.Get(ctx, v.Record)
}
