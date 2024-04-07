package dto4logist

import (
	"fmt"
	"github.com/strongo/validation"
	"strings"
)

// SetOrderCounterparty is a request to set a counterparty for an order
type SetOrderCounterparty struct {
	ContactID string `json:"contactID"`
	Role      string `json:"role"`
	RefNumber string `json:"refNumber"`
}

// Validate returns error if request is invalid
func (v SetOrderCounterparty) Validate() error {
	if strings.TrimSpace(v.ContactID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("contactID")
	}
	if strings.TrimSpace(v.Role) == "" {
		return validation.NewErrRequestIsMissingRequiredField("role")
	}
	return nil
}

// SetOrderCounterpartiesRequest is a request to set counterparties for an order
type SetOrderCounterpartiesRequest struct {
	OrderRequest
	Counterparties []SetOrderCounterparty `json:"counterparties"`
}

// Validate returns error if request is invalid
func (v SetOrderCounterpartiesRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	if len(v.Counterparties) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("counterparties")
	}
	for i, counterparty := range v.Counterparties {
		if err := counterparty.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("counterparties[%v]", i), err.Error())
		}
	}
	return nil
}
