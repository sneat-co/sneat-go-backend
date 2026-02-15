package dto4calendarium

import (
	"strconv"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dbo4calendarium"
	"github.com/strongo/validation"
)

// HappeningPricesRequest adds prices to happening
type HappeningPricesRequest struct {
	HappeningRequest
	dbo4calendarium.WithHappeningPrices
}

// Validate returns error if not valid
func (v HappeningPricesRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if len(v.Prices) == 0 {
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
				return validation.NewErrBadRecordFieldValue("priceIDs", "duplicate price ContactID: "+id)
			}
		}
	}
	return nil
}
