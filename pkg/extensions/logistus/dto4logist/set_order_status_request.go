package dto4logist

import (
	"strings"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/strongo/validation"
)

// SetOrderStatusRequest is a request to set status of an order
type SetOrderStatusRequest struct {
	OrderRequest
	Status dbo4logist.OrderStatus `json:"status"`
}

// Validate returns error if request is invalid
func (v SetOrderStatusRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(string(v.Status)) == "" {
		return validation.NewErrRequestIsMissingRequiredField("status")
	}
	return nil
}
