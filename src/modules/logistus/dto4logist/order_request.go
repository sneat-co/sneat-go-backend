package dto4logist

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/strongo/validation"
	"strings"
)

// NewOrderRequest creates new OrderRequest
func NewOrderRequest(spaceID coretypes.SpaceID, orderID string) OrderRequest {
	return OrderRequest{
		SpaceRequest: dto4spaceus.NewSpaceRequest(spaceID),
		OrderID:      orderID,
	}
}

// OrderRequest is a request regards an order that refers to a team with SpaceRequest and points to a specific order by OrderID
type OrderRequest struct {
	dto4spaceus.SpaceRequest
	OrderID string `json:"orderID"`
}

// Validate returns error if request is invalid
func (v OrderRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.OrderID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("orderID")
	}
	return nil
}
