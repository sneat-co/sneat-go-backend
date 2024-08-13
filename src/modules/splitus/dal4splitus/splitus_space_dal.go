package dal4splitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
)

func GetSplitusSpace(ctx context.Context, tx dal.ReadTransaction, splitusSpace models4splitus.SplitusSpaceEntry) (err error) {
	return tx.Get(ctx, splitusSpace.Record)
}

func SaveSplitusSpace(ctx context.Context, tx dal.ReadwriteTransaction, space models4splitus.SplitusSpaceEntry) error {
	return tx.Set(ctx, space.Record)
}
