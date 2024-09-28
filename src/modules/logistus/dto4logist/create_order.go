package dto4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/strongo/validation"
)

// CreateOrderRequest is a request to create an order
type CreateOrderRequest struct {
	dto4spaceus.SpaceRequest
	Order              dbo4logist.OrderBase `json:"order"`
	NumberOfContainers map[string]int       `json:"numberOfContainers"`
}

// Validate returns error if request is invalid
func (v CreateOrderRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}
	if err := v.Order.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}
	for k, v := range v.NumberOfContainers {
		if v < 0 {
			return validation.NewErrBadRequestFieldValue("numberOfContainers."+k, "should be >= 0")
		}
	}
	return nil
}

// CreateOrderResponse is a response to create an order request
type CreateOrderResponse struct {
	Order *dbo4logist.OrderBrief `json:"order"`
}
