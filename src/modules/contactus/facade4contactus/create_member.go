package facade4contactus

import (
	"context"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/checks4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// CreateMember adds members to a team
func CreateMember(
	ctx context.Context,
	user facade.User,
	request dal4contactus.CreateMemberRequest,
) (
	response dto4contactus.CreateContactResponse,
	err error,
) {
	if err = request.Validate(); err != nil {
		return response, fmt.Errorf("invalid CreateMemberRequest: %w", err)
	}
	createContactRequest := dto4contactus.CreateContactRequest{
		SpaceRequest: request.SpaceRequest,
		WithRelated:  request.WithRelated,
		Status:       request.Status,
		Type:         briefs4contactus.ContactTypePerson,
		Person:       &request.CreatePersonRequest,
	}
	if !checks4contactus.IsSpaceMember(request.Roles) {
		createContactRequest.Roles = append(createContactRequest.Roles, const4contactus.SpaceMemberRoleMember)
	}
	if response, err = CreateContact(ctx, user, false, createContactRequest); err != nil {
		return response, err
	}
	if response.Data == nil {
		return response, fmt.Errorf("CreateContact returned nil response data")
	}
	if !checks4contactus.IsSpaceMember(response.Data.Roles) {
		err = fmt.Errorf("created contact does not have team member role")
		return response, err
	}
	return response, err
}
