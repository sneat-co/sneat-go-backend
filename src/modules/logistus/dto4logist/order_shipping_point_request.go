package dto4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
	"github.com/strongo/validation"
)

// OrderShippingPointRequest represents a request that refers to an order shipping point.
type OrderShippingPointRequest struct {
	OrderRequest
	ShippingPointID string `json:"shippingPointID"`
}

// Validate returns an error if the request is invalid.
func (v OrderShippingPointRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("OrderRequest", err.Error())
	}
	if err := dbo4spaceus.ValidateShippingPointID(v.ShippingPointID); err != nil {
		return validation.NewErrBadRequestFieldValue("shippingPointID", err.Error())
	}
	return nil
}
