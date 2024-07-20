package dal4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

var _ facade.Request = (*CreateMemberRequest)(nil)

// CreateMemberRequest request is similar to dto4contactus.CreateContactRequest but has less fields
type CreateMemberRequest struct {
	dto4spaceus.SpaceRequest
	dto4contactus.CreatePersonRequest
	dbo4linkage.WithRelated
	Message string `json:"message"`
	//RelatedTo *dbo4linkage.RelationshipItemRolesCommand `json:"relatedTo,omitempty"`
}

// Validate validates request
func (v *CreateMemberRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
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
