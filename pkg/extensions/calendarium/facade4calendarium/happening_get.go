package facade4calendarium

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// GetByID returns RecurringHappeningDto record
func GetByID(ctx context.Context, getter dal.ReadSession, spaceID coretypes.SpaceID, happeningID string, dto dbo4calendarium.HappeningDbo) (record dal.Record, err error) {
	record = dal.NewRecordWithData(dbo4calendarium.NewHappeningKey(spaceID, happeningID), dto)
	return record, getter.Get(ctx, record)
}

// GetForUpdate returns SpaceIDs record in transaction
func GetForUpdate(ctx context.Context, tx dal.ReadwriteTransaction, spaceID coretypes.SpaceID, happeningID string, dto dbo4calendarium.HappeningDbo) (record dal.Record, err error) {
	return GetByID(ctx, tx, spaceID, happeningID, dto)
}
