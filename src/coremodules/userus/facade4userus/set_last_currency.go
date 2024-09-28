package facade4userus

import (
	"context"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func SetLastCurrency(ctx context.Context, userCtx facade.UserContext, currencyCode money.CurrencyCode) (err error) {
	return dal4userus.RunUserWorker(ctx, userCtx, true, func(ctx context.Context, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) (err error) {
		return txSetLastCurrency(ctx, tx, userWorkerParams, currencyCode)
	})
}

func txSetLastCurrency(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams, currencyCode money.CurrencyCode) (err error) {
	params.UserUpdates, err = params.User.Data.WithLastCurrencies.SetLastCurrency(currencyCode)
	return
}
