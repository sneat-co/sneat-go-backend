package dbo4budgetus

import (
	"github.com/crediterra/money"
	"testing"
)

func TestObligation_Validate(t *testing.T) {
	tests := []struct {
		name       string
		obligation Obligation
		wantErr    bool
	}{
		{"valid_collection", Obligation{
			Direction: ObligationDirectionCollection,
			Amount:    money.Amount{Value: 100, Currency: "USD"},
		}, false},
		{"valid_disbursement", Obligation{
			Direction: ObligationDirectionDisbursements,
			Amount:    money.Amount{Value: 100, Currency: "USD"},
		}, false},
		{"missing_direction", Obligation{
			Amount: money.Amount{Value: 100, Currency: "USD"},
		}, true},
		{"invalid_direction", Obligation{
			Direction: "invalid",
			Amount:    money.Amount{Value: 100, Currency: "USD"},
		}, true},
		{"invalid_amount", Obligation{
			Direction: ObligationDirectionCollection,
			Amount:    money.Amount{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.obligation.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Obligation.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
