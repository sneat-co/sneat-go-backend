package facade4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"slices"
)

func DeleteHappeningPrices(ctx context.Context, userCtx facade.UserContext, request dto4calendarium.DeleteHappeningPricesRequest) (err error) {
	var deleteHappeningPricesWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return deleteHappeningPricesTx(ctx, tx, userCtx, params, request)
	}
	return dal4calendarium.RunHappeningSpaceWorker(ctx, userCtx, request.HappeningRequest, deleteHappeningPricesWorker)
}

func deleteHappeningPricesTx(ctx context.Context, tx dal.ReadwriteTransaction, _ facade.UserContext, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.DeleteHappeningPricesRequest) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	happeningDbo := params.Happening.Data

	prices := make([]*dbo4calendarium.HappeningPrice, 0, len(happeningDbo.Prices))

	for _, price := range happeningDbo.Prices {
		if slices.Contains(request.PriceIDs, price.ID) {
			continue
		}
		prices = append(prices, price)
	}
	if len(prices) < len(happeningDbo.Prices) {
		happeningDbo.Prices = prices
		params.HappeningUpdates = append(params.HappeningUpdates, dal.Update{
			Field: "prices",
			Value: happeningDbo.Prices,
		})
		params.Happening.Record.MarkAsChanged()
	}
	return nil
}
