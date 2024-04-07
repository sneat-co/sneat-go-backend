package dto4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/strongo/validation"
	"strings"
)

// DeleteOrderCounterpartyRequest is a request to delete a counterparty from an order
type DeleteOrderCounterpartyRequest struct {
	OrderRequest
	Role      models4logist.CounterpartyRole `json:"role"`
	ContactID string                         `json:"contactID"`
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
