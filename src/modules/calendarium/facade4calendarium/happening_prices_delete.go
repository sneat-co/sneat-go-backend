package facade4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

func DeleteHappeningPrices(ctx context.Context, user facade.User, request dto4calendarium.DeleteHappeningPricesRequest) (err error) {
	var deleteHappeningPricesWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return deleteHappeningPricesTx(ctx, tx, user, params, request)
	}
	return dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, deleteHappeningPricesWorker)
}

func deleteHappeningPricesTx(ctx context.Context, tx dal.ReadwriteTransaction, _ facade.User, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.DeleteHappeningPricesRequest) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	happeningDbo := params.Happening.Dbo

	prices := make([]*models4calendarium.HappeningPrice, 0, len(happeningDbo.Prices))

	for _, price := range happeningDbo.Prices {
		if slice.Contains(request.PriceIDs, price.ID) {
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
