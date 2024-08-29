package dal4listus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
)

// GetListByID returns List record
func GetListByID(ctx context.Context, getter dal.ReadSession, list ListEntry) error {
	return getter.Get(ctx, list.Record)
}

// GetListForUpdate returns List record in read-write transaction
func GetListForUpdate(ctx context.Context, tx dal.ReadwriteTransaction, list ListEntry) error {
	return GetListByID(ctx, tx, list)
}
