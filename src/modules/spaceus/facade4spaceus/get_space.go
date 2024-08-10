package facade4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// GetSpace loads team record
func GetSpace(ctx context.Context, userContext facade.User, id string) (space dal4spaceus.SpaceEntry, err error) {
	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return space, err
	}
	space, err = GetSpaceByID(ctx, db, id)
	if err != nil || !space.Record.Exists() {
		return space, err
	}
	userID := userContext.GetID()
	var found bool
	for _, uid := range space.Data.UserIDs {
		if uid == userID {
			found = true
			break
		}
	}
	if !found {
		return space, fmt.Errorf("%w: you do not belong to the SpaceIDs", facade.ErrUnauthorized)
	}
	return space, err
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
