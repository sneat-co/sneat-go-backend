package facade4calendarius

import (
	"slices"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dal4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-core/facade"
)

func DeleteHappeningPrices(ctx facade.ContextWithUser, request dto4calendarius.DeleteHappeningPricesRequest) (err error) {
	var deleteHappeningPricesWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams) error {
		return deleteHappeningPricesTx(ctx, tx, params, request)
	}
	return dal4calendarius.RunHappeningSpaceWorker(ctx, request.HappeningRequest, deleteHappeningPricesWorker)
}

func deleteHappeningPricesTx(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *dal4calendarius.HappeningWorkerParams,
	request dto4calendarius.DeleteHappeningPricesRequest,
) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	happeningDbo := params.Happening.Data

	prices := make([]*dbo4calendarius.HappeningPrice, 0, len(happeningDbo.Prices))

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
