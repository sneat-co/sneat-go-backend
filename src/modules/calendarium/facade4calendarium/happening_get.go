package facade4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
)

// GetByID returns RecurringHappeningDto record
func GetByID(ctx context.Context, getter dal.ReadSession, teamID, happeningID string, dto dbo4calendarium.HappeningDbo) (record dal.Record, err error) {
	record = dal.NewRecordWithData(dbo4calendarium.NewHappeningKey(teamID, happeningID), dto)
	return record, getter.Get(ctx, record)
}

// GetForUpdate returns SpaceIDs record in transaction
func GetForUpdate(ctx context.Context, tx dal.ReadwriteTransaction, teamID, happeningID string, dto dbo4calendarium.HappeningDbo) (record dal.Record, err error) {
	return GetByID(ctx, tx, teamID, happeningID, dto)
}
