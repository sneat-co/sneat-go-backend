package dbo4calendarius

import (
	"testing"

	"github.com/crediterra/money"
	"github.com/stretchr/testify/assert"
)

func validAmount() money.Amount {
	return money.NewAmount(money.CurrencyUSD, 100)
}

func TestTerm_ID(t *testing.T) {
	tests := []struct {
		name string
		term Term
		want string
	}{
		{"single", Term{Unit: TermUnitSingle, Length: 1}, "single1"},
		{"hour", Term{Unit: TermUnitHour, Length: 2}, "hour2"},
		{"month_zero", Term{Unit: TermUnitMonth, Length: 0}, "month0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.term.ID())
		})
	}
}

func TestTerm_String(t *testing.T) {
	tests := []struct {
		name string
		term Term
		want string
	}{
		{"single", Term{Unit: TermUnitSingle, Length: 5}, "single"},
		{"length_one", Term{Unit: TermUnitHour, Length: 1}, "1 hour"},
		{"length_many", Term{Unit: TermUnitDay, Length: 3}, "3 days"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.term.String())
		})
	}
}

func TestTerm_Validate(t *testing.T) {
	tests := []struct {
		name    string
		term    Term
		wantErr bool
	}{
		{"valid_single", Term{Unit: TermUnitSingle, Length: 1}, false},
		{"valid_hour", Term{Unit: TermUnitHour, Length: 2}, false},
		{"empty_unit", Term{Unit: "", Length: 1}, true},
		{"unknown_unit", Term{Unit: "decade", Length: 1}, true},
		{"zero_length", Term{Unit: TermUnitHour, Length: 0}, true},
		{"negative_length", Term{Unit: TermUnitHour, Length: -1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.term.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHappeningPrice_Validate(t *testing.T) {
	validTerm := Term{Unit: TermUnitHour, Length: 1}
	tests := []struct {
		name    string
		price   HappeningPrice
		wantErr bool
	}{
		{"valid", HappeningPrice{ID: "p1", Term: validTerm, Amount: validAmount()}, false},
		{"valid_empty_id", HappeningPrice{Term: validTerm, Amount: validAmount()}, false},
		{"valid_with_expense_qty", HappeningPrice{ID: "p1", Term: validTerm, Amount: validAmount(), ExpenseQuantity: 2}, false},
		{"id_star", HappeningPrice{ID: "*", Term: validTerm, Amount: validAmount()}, true},
		{"invalid_term", HappeningPrice{ID: "p1", Term: Term{}, Amount: validAmount()}, true},
		{"invalid_amount", HappeningPrice{ID: "p1", Term: validTerm, Amount: money.Amount{Value: 100}}, true},
		{"negative_amount", HappeningPrice{ID: "p1", Term: validTerm, Amount: money.NewAmount(money.CurrencyUSD, -1)}, true},
		{"negative_expense_qty", HappeningPrice{ID: "p1", Term: validTerm, Amount: validAmount(), ExpenseQuantity: -1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.price.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWithHappeningPrices_GetPriceByID(t *testing.T) {
	p1 := &HappeningPrice{ID: "p1", Term: Term{Unit: TermUnitHour, Length: 1}, Amount: validAmount()}
	p2 := &HappeningPrice{ID: "p2", Term: Term{Unit: TermUnitDay, Length: 1}, Amount: validAmount()}
	v := WithHappeningPrices{Prices: []*HappeningPrice{p1, p2}}

	assert.Same(t, p1, v.GetPriceByID("p1"))
	assert.Same(t, p2, v.GetPriceByID("p2"))
	assert.Nil(t, v.GetPriceByID("missing"))
	assert.Nil(t, WithHappeningPrices{}.GetPriceByID("p1"))
}

func TestWithHappeningPrices_Validate(t *testing.T) {
	validTerm := Term{Unit: TermUnitHour, Length: 1}
	tests := []struct {
		name    string
		v       WithHappeningPrices
		wantErr bool
	}{
		{"empty", WithHappeningPrices{}, false},
		{
			name: "valid",
			v: WithHappeningPrices{Prices: []*HappeningPrice{
				{ID: "p1", Term: Term{Unit: TermUnitHour, Length: 1}, Amount: validAmount()},
				{ID: "p2", Term: Term{Unit: TermUnitDay, Length: 1}, Amount: validAmount()},
			}},
			wantErr: false,
		},
		{
			name:    "nil_price",
			v:       WithHappeningPrices{Prices: []*HappeningPrice{nil}},
			wantErr: true,
		},
		{
			name: "invalid_price",
			v: WithHappeningPrices{Prices: []*HappeningPrice{
				{ID: "p1", Term: Term{}, Amount: validAmount()},
			}},
			wantErr: true,
		},
		{
			// Validate flags a price whose ID equals another price's Term.ID().
			name: "id_collides_with_term_id",
			v: WithHappeningPrices{Prices: []*HappeningPrice{
				{ID: validTerm.ID(), Term: Term{Unit: TermUnitDay, Length: 1}, Amount: validAmount()},
				{ID: "p2", Term: validTerm, Amount: validAmount()},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
