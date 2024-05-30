package dal4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

var _ facade.Request = (*CreateMemberRequest)(nil)

// CreateMemberRequest request is similar to dto4contactus.CreateContactRequest but has less fields
type CreateMemberRequest struct {
	dto4teamus.TeamRequest
	dto4contactus.CreatePersonRequest
	models4linkage.WithRelated
	Message string `json:"message"`
	//RelatedTo *models4linkage.RelationshipRolesCommand `json:"relatedTo,omitempty"`
}

// Validate validates request
func (v *CreateMemberRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if err := v.CreatePersonRequest.Validate(); err != nil {
		return err
	}
	if err := v.WithRelated.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}
	return nil
}
