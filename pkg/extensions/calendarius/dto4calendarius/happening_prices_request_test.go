package dto4calendarius

import (
	"testing"

	"github.com/crediterra/money"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/stretchr/testify/assert"
)

func validHappeningPrice() *dbo4calendarius.HappeningPrice {
	return &dbo4calendarius.HappeningPrice{
		ID:     "p1",
		Term:   dbo4calendarius.Term{Unit: dbo4calendarius.TermUnitHour, Length: 1},
		Amount: money.NewAmount(money.CurrencyUSD, 100),
	}
}

func TestHappeningPricesRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     HappeningPricesRequest
		wantErr bool
	}{
		{"valid", HappeningPricesRequest{
			HappeningRequest: validHappeningRequest(),
			WithHappeningPrices: dbo4calendarius.WithHappeningPrices{
				Prices: []*dbo4calendarius.HappeningPrice{validHappeningPrice()},
			},
		}, false},
		{"invalid_happening_request", HappeningPricesRequest{
			HappeningRequest: HappeningRequest{}, // missing space + happeningID
			WithHappeningPrices: dbo4calendarius.WithHappeningPrices{
				Prices: []*dbo4calendarius.HappeningPrice{validHappeningPrice()},
			},
		}, true},
		{"empty_prices", HappeningPricesRequest{
			HappeningRequest: validHappeningRequest(),
		}, true},
		{"invalid_price", HappeningPricesRequest{
			HappeningRequest: validHappeningRequest(),
			WithHappeningPrices: dbo4calendarius.WithHappeningPrices{
				Prices: []*dbo4calendarius.HappeningPrice{{ID: "p1", Term: dbo4calendarius.Term{}}},
			},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteHappeningPricesRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     DeleteHappeningPricesRequest
		wantErr bool
	}{
		{"valid_single", DeleteHappeningPricesRequest{
			HappeningRequest: validHappeningRequest(),
			PriceIDs:         []string{"p1"},
		}, false},
		{"distinct_ids_treated_as_dup", DeleteHappeningPricesRequest{
			HappeningRequest: validHappeningRequest(),
			PriceIDs:         []string{"p1", "p2"},
		}, true},
		{"invalid_happening_request", DeleteHappeningPricesRequest{
			HappeningRequest: HappeningRequest{},
			PriceIDs:         []string{"p1"},
		}, true},
		{"empty_price_ids", DeleteHappeningPricesRequest{
			HappeningRequest: validHappeningRequest(),
		}, true},
		{"empty_string_id", DeleteHappeningPricesRequest{
			HappeningRequest: validHappeningRequest(),
			PriceIDs:         []string{""},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
