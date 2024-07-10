package facade4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
)

// CreateSpaceItemRequest DTO
type CreateSpaceItemRequest struct {
	dto4teamus.SpaceRequest
	ContactID string `json:"contactID,omitempty"`
	MemberID  string `json:"memberID,omitempty"`
}

// Validate returns error if not valid
func (v CreateSpaceItemRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	return nil
}
