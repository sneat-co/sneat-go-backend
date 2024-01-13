package facade4teamus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// GetTeam loads team record
func GetTeam(ctx context.Context, userContext facade.User, id string) (team dal4teamus.TeamContext, err error) {
	db := facade.GetDatabase(ctx)
	var record dal.Record
	team, err = GetTeamByID(ctx, db, id)
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
		return team, fmt.Errorf("%w: you do not belong to the TeamIDs", facade.ErrUnauthorized)
	}
	return team, err
}

// GetTeamByID return TeamIDs record
func GetTeamByID(ctx context.Context, getter dal.ReadSession, id string) (team dal4teamus.TeamContext, err error) {
	team = dal4teamus.NewTeamContext(id)
	return team, getter.Get(ctx, team.Record)
}

// TxGetTeamByID returns TeamIDs record in transaction
func TxGetTeamByID(ctx context.Context, tx dal.ReadwriteTransaction, id string) (team dal4teamus.TeamContext, err error) {
	return GetTeamByID(ctx, tx, id)
}
