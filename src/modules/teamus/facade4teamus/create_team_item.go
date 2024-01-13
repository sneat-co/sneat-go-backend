package facade4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
)

// CreateTeamItemRequest DTO
type CreateTeamItemRequest struct {
	dto4teamus.TeamRequest
	ContactID string `json:"contactID,omitempty"`
	MemberID  string `json:"memberID,omitempty"`
}

// Validate returns error if not valid
func (v CreateTeamItemRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	return nil
}
