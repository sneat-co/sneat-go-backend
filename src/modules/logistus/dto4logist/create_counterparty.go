package dto4logist

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
)

// CreateCounterpartyRequest defines a request to create a new counterparty
type CreateCounterpartyRequest struct {
	dto4spaceus.SpaceRequest
	with.RolesField
	Company dto4contactus.CreateCompanyRequest `json:"company"`
}

// Validate returns error if request is invalid
func (v CreateCounterpartyRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if err := v.Company.Validate(); err != nil {
		return err
	}
	if err := v.RolesField.Validate(); err != nil {
		return fmt.Errorf("%w: %v", facade.ErrBadRequest, err.Error())
	}
	return nil
}
