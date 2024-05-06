package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/random"
)

func AddHappeningPrices(ctx context.Context, user facade.User, request dto4calendarium.HappeningPricesRequest) (err error) {
	var addHappeningPricesWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return addHappeningPricesTx(ctx, tx, user, params, request)
	}
	return dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, addHappeningPricesWorker)
}

func addHappeningPricesTx(ctx context.Context, tx dal.ReadwriteTransaction, _ facade.User, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.HappeningPricesRequest) (err error) {
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
		termID := price.Term.ID()
		price.ID = termID // Ignore ID passed in request from client
		const maxAttempts = 10
		for k := 0; k < maxAttempts+1; k++ {
			if k == maxAttempts {
				return fmt.Errorf("too many attempts to generate unique price ID for term %v", termID)
			}
			if price.ID != "" && happeningDbo.GetPriceByID(price.ID) == nil {
				happeningDbo.Prices = append(happeningDbo.Prices, price)
				break
			}
			price.ID = termID + random.ID(k+1)
		}
		params.Happening.Record.MarkAsChanged()
	}

	if params.Happening.Record.HasChanged() {
		params.HappeningUpdates = append(params.HappeningUpdates, dal.Update{
			Field: "prices",
			Value: happeningDbo.Prices,
		})
	}

	return nil
}
