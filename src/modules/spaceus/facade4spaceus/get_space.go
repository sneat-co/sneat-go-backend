package facade4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// GetSpace loads team record
func GetSpace(ctx context.Context, userCtx facade.UserContext, id string) (space dbo4spaceus.SpaceEntry, err error) {
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return space, err
	}
	space, err = GetSpaceByID(ctx, db, id)
	if err != nil || !space.Record.Exists() {
		return space, err
	}
	userID := userCtx.GetUserID()
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
func GetSpaceByID(ctx context.Context, getter dal.ReadSession, id string) (space dbo4spaceus.SpaceEntry, err error) {
	space = dbo4spaceus.NewSpaceEntry(id)
	return space, getter.Get(ctx, space.Record)
}

// TxGetSpaceByID returns SpaceIDs record in transaction
func TxGetSpaceByID(ctx context.Context, tx dal.ReadwriteTransaction, id string) (team dbo4spaceus.SpaceEntry, err error) {
	return GetSpaceByID(ctx, tx, id)
}
