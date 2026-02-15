package facade4calendarium

import (
	"slices"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
)

func DeleteHappeningPrices(ctx facade.ContextWithUser, request dto4calendarium.DeleteHappeningPricesRequest) (err error) {
	var deleteHappeningPricesWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return deleteHappeningPricesTx(ctx, tx, params, request)
	}
	return dal4calendarium.RunHappeningSpaceWorker(ctx, request.HappeningRequest, deleteHappeningPricesWorker)
}

func deleteHappeningPricesTx(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *dal4calendarium.HappeningWorkerParams,
	request dto4calendarium.DeleteHappeningPricesRequest,
) (err error) {
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
		params.HappeningUpdates = append(params.HappeningUpdates,
			update.ByFieldName("prices", happeningDbo.Prices))
		params.Happening.Record.MarkAsChanged()
	}
	return nil
}
