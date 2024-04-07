package dto4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/strongo/validation"
	"strings"
)

// SetOrderStatusRequest is a request to set status of an order
type SetOrderStatusRequest struct {
	OrderRequest
	Status models4logist.OrderStatus `json:"status"`
}

// Validate returns error if request is invalid
func (v SetOrderStatusRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(string(v.Status)) == "" {
		return validation.NewErrRequestIsMissingRequiredField("status")
	}
	return nil
}
