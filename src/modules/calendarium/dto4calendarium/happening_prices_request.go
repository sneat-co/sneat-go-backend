package dto4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/strongo/validation"
	"strconv"
)

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
	if len(v.WithHappeningPrices.Prices) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("prices")
	}
	if err := v.WithHappeningPrices.Validate(); err != nil {
		return err
	}
	return nil
}

// DeleteHappeningPricesRequest adds prices to happening
type DeleteHappeningPricesRequest struct {
	HappeningRequest
	PriceIDs []string `json:"priceIDs"`
}

func (v DeleteHappeningPricesRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if len(v.PriceIDs) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("priceIDs")
	}
	for i, id := range v.PriceIDs {
		if id == "" {
			return validation.NewErrBadRecordFieldValue("priceIDs["+strconv.Itoa(i)+"]", "empty value")
		}
		for j, id2 := range v.PriceIDs {
			if i != j && id != id2 {
				return validation.NewErrBadRecordFieldValue("priceIDs", "duplicate price ID: "+id)
			}
		}
	}
	return nil
}
