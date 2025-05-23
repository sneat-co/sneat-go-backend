package facade4calendarium

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/random"
)

func SetHappeningPrices(ctx facade.ContextWithUser, request dto4calendarium.HappeningPricesRequest) (err error) {
	var setHappeningPricesWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return setHappeningPricesTx(ctx, tx, params, request)
	}
	return dal4calendarium.RunHappeningSpaceWorker(ctx, request.HappeningRequest, setHappeningPricesWorker)
}

func setHappeningPricesTx(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.HappeningPricesRequest) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	if !params.Happening.Record.Exists() {
		return fmt.Errorf("happening not found by key=%s: %v", params.Happening.Key.String(), params.Happening.Record.Error())
	}
	happeningDbo := params.Happening.Data

requestPrices:
	for _, requestPrice := range request.Prices {

		// Check if we are updating existing price
		for _, dbPrice := range happeningDbo.Prices {
			if requestPrice.ID != "" && requestPrice.ID == dbPrice.ID {
				if dbPrice.Amount != requestPrice.Amount {
					dbPrice.Amount = requestPrice.Amount
					params.Happening.Record.MarkAsChanged()
				}
				if dbPrice.ExpenseQuantity != requestPrice.ExpenseQuantity {
					dbPrice.ExpenseQuantity = requestPrice.ExpenseQuantity
					params.Happening.Record.MarkAsChanged()
				}
				continue requestPrices
			}
		}

		termID := requestPrice.Term.ID()
		requestPrice.ID = termID // Ignore ContactID passed in request from client
		const maxAttempts = 10
		for k := 0; k < maxAttempts+1; k++ {
			if k == maxAttempts {
				return fmt.Errorf("too many attempts to generate unique requestPrice ContactID for term %v", termID)
			}
			if requestPrice.ID != "" && happeningDbo.GetPriceByID(requestPrice.ID) == nil {
				happeningDbo.Prices = append(happeningDbo.Prices, requestPrice)
				break
			}
			requestPrice.ID = termID + random.ID(k+1)
		}
		params.Happening.Record.MarkAsChanged()
	}

	if params.Happening.Record.HasChanged() {
		params.HappeningUpdates = append(params.HappeningUpdates,
			update.ByFieldName("prices", happeningDbo.Prices))
	}

	return nil
}
