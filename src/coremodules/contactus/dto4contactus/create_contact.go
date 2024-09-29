package dto4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
)

// CreateContactRequest DTO
type CreateContactRequest struct {
	dto4spaceus.SpaceRequest
	with.RolesField
	ParentContactID string                       `json:"parentContactID,omitempty"`
	Type            briefs4contactus.ContactType `json:"type"`

	// Duplicate also in CreatePersonRequest throw briefs4contactus.ContactBase,
	// but not in CreateCompanyRequest & CreateLocationRequest
	Status string `json:"status"`

	Person   *CreatePersonRequest       `json:"person,omitempty"`
	Company  *CreateCompanyRequest      `json:"company,omitempty"`
	Location *CreateLocationRequest     `json:"location,omitempty"`
	Basic    *CreateBasicContactRequest `json:"basic,omitempty"`

	dbo4linkage.WithRelated

	// Used for situation when we want a hard-coded contact number
	// (e.g. a self-contact for a company team).
	// Can not be used from client side
	ContactID string `json:"-"`
}

func (v CreateContactRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	switch v.Status {
	case "":
		return validation.NewErrRequestIsMissingRequiredField("status")
	case "active", "draft":
		// OK
	default:
		return validation.NewErrBadRequestFieldValue("status", "allowed values are 'active' and 'draft', got: "+v.Status)
	}
	switch v.Type {
	case "":
		return validation.NewErrRequestIsMissingRequiredField("type")
	case "person":
		if v.Person == nil {
			return validation.NewErrRequestIsMissingRequiredField("person")
		}
		if err := v.Person.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("person", fmt.Sprintf("contact type is set to 'person', but the `person` field is invalid: %v", err))
		}
	case "company":
		if v.Company == nil {
			return validation.NewErrRequestIsMissingRequiredField("company")
		}
		if err := v.Company.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("company", fmt.Sprintf("contact type is set to 'company', but the `company` field is invalid: %v", err))
		}
	case "location":
		if v.Location == nil {
			return validation.NewErrRequestIsMissingRequiredField("location")
		}
		if err := v.Location.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("location", fmt.Sprintf("contact type is set to 'location', but the `location` field is invalid: %v", err))
		}
		if v.ParentContactID == "" {
			return validation.NewErrRequestIsMissingRequiredField("parentContactID")
		}
	default:
		if v.Basic == nil {
			return validation.NewErrRequestIsMissingRequiredField("basic")
		}
		if err := v.Basic.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("company", err.Error())
		}
	}
	if v.Person != nil && v.Type != "person" {
		return validation.NewErrBadRequestFieldValue("person", "the `person` field is not nil, but contact type is set to 'person'")
	}
	if v.Company != nil && v.Type != "company" {
		return validation.NewErrBadRequestFieldValue("company", "the `company` field is not nil, but contact type is set to 'company'")
	}
	if v.Location != nil && v.Type != "location" {
		return validation.NewErrBadRequestFieldValue("location", "the `location` field is not nil, but contact type is set to 'location'")
	}
	if err := v.RolesField.Validate(); err != nil {
		return fmt.Errorf("%w: %v", facade.ErrBadRequest, err)
	}
	if v.Person != nil && v.Person.Status != v.Status {
		return validation.NewErrBadRecordFieldValue("status",
			fmt.Sprintf("does not match to person.status: %s != %s", v.Status, v.Person.Status))
	}
	if err := v.WithRelated.Validate(); err != nil {
		return err
	}
	return nil
}

// CreateContactResponse DTO
type CreateContactResponse = dbmodels.DtoWithID[*dbo4contactus.ContactDbo]
