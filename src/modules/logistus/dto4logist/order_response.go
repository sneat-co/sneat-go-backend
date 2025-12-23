package dto4logist

import (
	"errors"

	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
)

// OrderResponse is a response to an order modification request.
type OrderResponse struct {
	OrderDto *dbo4logist.OrderDbo `json:"order"`
}

// Validate returns an error if the response is invalid.
func (v OrderResponse) Validate() error {
	if v.OrderDto == nil {
		return errors.New("response is missing required field: order")
	}
	return nil
}
