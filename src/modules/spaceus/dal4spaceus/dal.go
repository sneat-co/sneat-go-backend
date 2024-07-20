package dal4spaceus

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

func txUpdateSpace(ctx context.Context, tx dal.ReadwriteTransaction, timestamp time.Time, space SpaceEntry, data []dal.Update, opts ...dal.Precondition) error {
	if err := space.Data.Validate(); err != nil {
		return fmt.Errorf("space record is not valid: %w", err)
	}
	space.Data.Version++
	data = append(data,
		dal.Update{Field: "v", Value: space.Data.Version},
		dal.Update{Field: "timestamp", Value: timestamp},
	)
	return txUpdate(ctx, tx, space.Key, data, opts...)
}

func txUpdateSpaceModule[D SpaceModuleDbo](ctx context.Context, tx dal.ReadwriteTransaction, _ time.Time, spaceModule record.DataWithID[string, D], data []dal.Update, opts ...dal.Precondition) error {
	if !spaceModule.Record.Exists() {
		return fmt.Errorf("an attempt to update a space module record that does not exist: %s", spaceModule.Key.String())
	}
	return txUpdate(ctx, tx, spaceModule.Key, data, opts...)
}
