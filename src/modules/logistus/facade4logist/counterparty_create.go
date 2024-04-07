package facade4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// CreateCounterparty creates a new counterparty
func CreateCounterparty(
	ctx context.Context,
	userContext facade.User,
	request dto4logist.CreateCounterpartyRequest,
) (
	response dto4contactus.CreateContactResponse, err error,
) {
	if len(request.Roles) == 0 {
		return response, validation.NewErrRequestIsMissingRequiredField("company.roles")
	}
	createContactRequest := dto4contactus.CreateContactRequest{
		TeamRequest: request.TeamRequest,
		Company:     &request.Company,
	}
	return facade4contactus.CreateContact(ctx, userContext, false, createContactRequest)
}
