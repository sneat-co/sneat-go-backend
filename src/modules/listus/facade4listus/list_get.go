package facade4listus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
)

// GetListByID returns List record
func GetListByID(ctx context.Context, getter dal.ReadSession, list dal4listus.ListEntry) error {
	return getter.Get(ctx, list.Record)
}

// GetListForUpdate returns List record in read-write transaction
func GetListForUpdate(ctx context.Context, tx dal.ReadwriteTransaction, list dal4listus.ListEntry) error {
	return GetListByID(ctx, tx, list)
}
