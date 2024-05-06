package dto4calendarium

import "github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"

// HappeningPricesRequest adds prices to happening
type HappeningPricesRequest struct {
	HappeningRequest
	models4calendarium.WithHappeningPrices
}

// Validate returns error if not valid
func (v HappeningPricesRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if err := v.WithHappeningPrices.Validate(); err != nil {
		return err
	}
	return nil
}
