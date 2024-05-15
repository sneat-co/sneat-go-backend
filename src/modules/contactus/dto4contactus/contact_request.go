package dto4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/strongo/validation"
	"strings"
)

// ContactRequest defines a request for a single contact
type ContactRequest struct {
	dto4teamus.TeamRequest
	ContactID string `json:"contactID"`
}

// Validate returns error if request is invalid
func (v ContactRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.ContactID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("contactID")
	}
	return nil
}

// ContactsRequest defines a request for a single contact
type ContactsRequest struct {
	dto4teamus.TeamRequest
	ContactIDs []string `json:"contactIDs"`
}

// Validate returns error if request is invalid
func (v ContactsRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if len(v.ContactIDs) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("contactIDs")
	}
	for i, contactID := range v.ContactIDs {
		if strings.TrimSpace(contactID) == "" {
			return validation.NewErrRequestIsMissingRequiredField(fmt.Sprintf("contactIDs[%d]", i))
		}
	}
	return nil
}
