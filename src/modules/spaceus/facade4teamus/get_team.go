package facade4teamus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// GetSpace loads team record
func GetSpace(ctx context.Context, userContext facade.User, id string) (team dal4spaceus.SpaceEntry, err error) {
	db := facade.GetDatabase(ctx)
	var record dal.Record
	team, err = GetSpaceByID(ctx, db, id)
	if err != nil || !record.Exists() {
		return team, err
	}
	userID := userContext.GetID()
	var found bool
	for _, uid := range team.Data.UserIDs {
		if uid == userID {
			found = true
			break
		}
	}
	if !found {
		return team, fmt.Errorf("%w: you do not belong to the SpaceIDs", facade.ErrUnauthorized)
	}
	return team, err
}

// GetSpaceByID return SpaceIDs record
func GetSpaceByID(ctx context.Context, getter dal.ReadSession, id string) (team dal4spaceus.SpaceEntry, err error) {
	team = dal4spaceus.NewSpaceEntry(id)
	return team, getter.Get(ctx, team.Record)
}

// TxGetSpaceByID returns SpaceIDs record in transaction
func TxGetSpaceByID(ctx context.Context, tx dal.ReadwriteTransaction, id string) (team dal4spaceus.SpaceEntry, err error) {
	return GetSpaceByID(ctx, tx, id)
}
