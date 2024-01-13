package facade4userus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
)

var runReadwriteTransaction = func(ctx context.Context, f dal.RWTxWorker, options ...dal.TransactionOption) error {
	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, f, options...)
}
