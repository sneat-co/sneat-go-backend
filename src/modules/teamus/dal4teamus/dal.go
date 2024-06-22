package dal4teamus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-core/facade/db"
	"time"
)

var txUpdate = func(ctx context.Context, tx dal.ReadwriteTransaction, key *dal.Key, data []dal.Update, opts ...dal.Precondition) error {
	return db.TxUpdate(ctx, tx, key, data, opts...)
}

func txUpdateTeam(ctx context.Context, tx dal.ReadwriteTransaction, timestamp time.Time, team TeamEntry, data []dal.Update, opts ...dal.Precondition) error {
	if err := team.Data.Validate(); err != nil {
		return fmt.Errorf("team record is not valid: %w", err)
	}
	team.Data.Version++
	data = append(data,
		dal.Update{Field: "v", Value: team.Data.Version},
		dal.Update{Field: "timestamp", Value: timestamp},
	)
	return txUpdate(ctx, tx, team.Key, data, opts...)
}

func txUpdateTeamModule[D TeamModuleDbo](ctx context.Context, tx dal.ReadwriteTransaction, _ time.Time, teamModule record.DataWithID[string, D], data []dal.Update, opts ...dal.Precondition) error {
	if !teamModule.Record.Exists() {
		return fmt.Errorf("an attempt to update a team module record that does not exist: %s", teamModule.Key.String())
	}
	return txUpdate(ctx, tx, teamModule.Key, data, opts...)
}
