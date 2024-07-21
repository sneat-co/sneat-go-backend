package facade4spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
)

// CreateSpaceItemRequest DTO
type CreateSpaceItemRequest struct {
	dto4spaceus.SpaceRequest
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
