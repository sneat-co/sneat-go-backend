package facade4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
)

func DeleteHappeningPrices(ctx context.Context, user facade.User, request dto4calendarium.HappeningPricesRequest) (err error) {
	var deleteHappeningPricesWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return deleteHappeningPricesTx(ctx, tx, user, params, request)
	}
	return dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, deleteHappeningPricesWorker)
}

func deleteHappeningPricesTx(ctx context.Context, tx dal.ReadwriteTransaction, _ facade.User, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.HappeningPricesRequest) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	happeningDbo := params.Happening.Dbo

	prices := make([]*models4calendarium.HappeningPrice, 0, len(happeningDbo.Prices))

requestPrices:
	for _, price := range request.Prices {
		for _, p := range happeningDbo.Prices {
			if p.Term == price.Term {
				continue requestPrices
			}
		}
		prices = append(happeningDbo.Prices, price)
	}
	if len(prices) < len(happeningDbo.Prices) {
		happeningDbo.Prices = prices
		params.Happening.Record.MarkAsChanged()
	}
	return nil
}
