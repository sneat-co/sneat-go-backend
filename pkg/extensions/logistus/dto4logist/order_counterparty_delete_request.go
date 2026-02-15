package dto4logist

import (
	"strings"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/strongo/validation"
)

// DeleteOrderCounterpartyRequest is a request to delete a counterparty from an order
type DeleteOrderCounterpartyRequest struct {
	OrderRequest
	Role      dbo4logist.CounterpartyRole `json:"role"`
	ContactID string                      `json:"contactID"`
}

// Validate returns error if request is invalid
func (v DeleteOrderCounterpartyRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.Role) == "" {
		return validation.NewErrRequestIsMissingRequiredField("role")
	}
	if strings.TrimSpace(v.ContactID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("contactID")
	}
	return nil
}
