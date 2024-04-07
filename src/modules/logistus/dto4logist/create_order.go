package dto4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/strongo/validation"
)

// CreateOrderRequest is a request to create an order
type CreateOrderRequest struct {
	dto4teamus.TeamRequest
	Order              models4logist.OrderBase `json:"order"`
	NumberOfContainers map[string]int          `json:"numberOfContainers"`
}

// Validate returns error if request is invalid
func (v CreateOrderRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
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
	Order *models4logist.OrderBrief `json:"order"`
}
