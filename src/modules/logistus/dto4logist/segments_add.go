package dto4logist

import (
	"fmt"
	"strings"

	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

// SegmentContainerData represents a container for a segment
type SegmentContainerData struct {
	ID string `json:"id"`
	dbo4logist.FreightPoint
}

// Validate returns nil if the SegmentContainerData is valid, otherwise returns the first error.
func (v SegmentContainerData) Validate() error {
	if v.ID == "" {
		return validation.NewErrRequestIsMissingRequiredField("id")
	}
	if v.ToLoad != nil {
		if err := v.ToLoad.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("toPick", err.Error())
		}
	}
	if v.ToUnload != nil {
		if err := v.ToUnload.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("toDrop", err.Error())
		}
	}
	return nil
}

// AddSegmentParty represents a request to add a segment party to an order
type AddSegmentParty struct {
	Counterparty dbo4logist.SegmentCounterparty `json:"counterparty"`
	RefNumber    string                         `json:"refNumber,omitempty"`
}

// Validate returns nil if the AddSegmentParty request is valid, otherwise returns the first error.
func (v AddSegmentParty) Validate() error {
	if err := v.Counterparty.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("counterparty", err.Error())
	}
	if strings.TrimSpace(v.RefNumber) != v.RefNumber {
		return validation.NewErrBadRecordFieldValue("refNumber", "must be trimmed")
	}
	return nil
}

// AddSegmentEndpoint represents a request to add a segment endpoint to an order
type AddSegmentEndpoint struct {
	AddSegmentParty
	Date string `json:"date,omitempty"`
}

// Validate	returns nil if the AddSegmentEndpoint request is valid, otherwise returns the first error.
func (v AddSegmentEndpoint) Validate() error {
	if err := v.AddSegmentParty.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("AddSegmentParty", err.Error())
	}
	if v.Date != "" {
		if _, err := validate.DateString(v.Date); err != nil {
			return validation.NewErrBadRecordFieldValue("date", err.Error())
		}
	}
	return nil
}

// AddSegmentsRequest represents a request to add segments to an order
type AddSegmentsRequest struct {
	OrderRequest
	From       AddSegmentEndpoint     `json:"from"`
	To         AddSegmentEndpoint     `json:"to"`
	By         *AddSegmentParty       `json:"by,omitempty"`
	Containers []SegmentContainerData `json:"containers"`
}

// Validate returns nil if the AddSegmentsRequest is valid, otherwise returns the first error.
func (v AddSegmentsRequest) Validate() error {
	if err := v.From.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("from", err.Error())
	}
	if err := v.To.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("to", err.Error())
	}
	if v.By != nil {
		if err := v.By.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("by", err.Error())
		}
		if v.By.Counterparty.Role != dbo4logist.CounterpartyRoleTrucker {
			return validation.NewErrBadRequestFieldValue("by.counterparty.role", fmt.Sprintf("expected to be `trucker`, got: %v", v.By.Counterparty.Role))
		}
	}
	if len(v.Containers) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("containers")
	}
	for i, container := range v.Containers {
		if err := container.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("containers[%v]", i), err.Error())
		}
		for j, c2 := range v.Containers {
			if j == i {
				continue
			}
			if c2.ID == container.ID {
				return validation.NewErrBadRequestFieldValue(fmt.Sprintf("containers[%v]", i), "duplicate container ContactID")
			}
		}
	}
	return nil
}
