package facade4calendarium

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
)

func AddHappeningPrices(ctx context.Context, user facade.User, request dto4calendarium.HappeningPricesRequest) (err error) {
	var addHappeningPricesWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return addHappeningPricesTx(ctx, tx, user, params, request)
	}
	return dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, addHappeningPricesWorker)
}

func addHappeningPricesTx(ctx context.Context, tx dal.ReadwriteTransaction, user facade.User, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.HappeningPricesRequest) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	happeningDbo := params.Happening.Dbo

requestPrices:
	for _, price := range request.Prices {
		for _, p := range happeningDbo.Prices {
			if p.Term == price.Term {
				if p.Amount != price.Amount {
					p.Amount = price.Amount
					params.Happening.Record.MarkAsChanged()
				}
				continue requestPrices
			}
		}
		happeningDbo.Prices = append(happeningDbo.Prices, price)
		params.Happening.Record.MarkAsChanged()
	}

	return nil
}
