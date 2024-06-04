package dbo4logist

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

// WithOrderContacts is a struct that contains contacts in OrderDto
type WithOrderContacts struct {
	Contacts []*OrderContact `json:"contacts" firestore:"contacts"`
}

// Validate validates contacts
func (v WithOrderContacts) Validate() error {
	if len(v.Contacts) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("contacts")
	}
	contactIDs := make([]string, len(v.Contacts))
	for i, contact := range v.Contacts {
		if err := contact.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("contacts[%d]{id=%s}", i, contact.ID),
				err.Error(),
			)
		}
		if j := slice.Index(contactIDs, contact.ID); j >= 0 {
			return validation.NewErrBadRecordFieldValue("contacts", fmt.Sprintf("duplicate contact ID at index %d & %d: ID=[%s]", j, i, contact.ID))
		}
		contactIDs = append(contactIDs, contact.ID)
		if contact.ParentID != "" {
			_, parent := v.GetContactByID(contact.ParentID)
			if parent == nil {
				return fmt.Errorf("parent contact not found in `contacts` by ID=[%s]", contact.ParentID)
			}
		}
		//if contact.Type == dbmodels.ContactTypeLocation && strings.TrimSpace(contact.Address.Lines) == "" {
		//	return validation.NewErrRecordIsMissingRequiredField("address.lines")
		//}
	}
	return nil
}

// Updates returns updates for WithCounterparties
func (v WithOrderContacts) Updates() []dal.Update {
	return []dal.Update{
		{Field: "contacts", Value: v.Contacts},
	}
}

// GetContactByID returns contact by ID
func (v WithOrderContacts) GetContactByID(id string) (int, *OrderContact) {
	for i, contact := range v.Contacts {
		if contact.ID == id {
			return i, contact
		}
	}
	return -1, nil
}

// MustGetContactByID must return contact by ID
func (v WithOrderContacts) MustGetContactByID(id string) *OrderContact {
	if strings.TrimSpace(id) == "" {
		panic("id is required parameter")
	}
	for _, contact := range v.Contacts {
		if contact.ID == id {
			return contact
		}
	}
	panic("contact not found by ID=" + id)
}

// GetContactByParentID returns first contact by parent ID
func (v WithOrderContacts) GetContactByParentID(parentID string) (int, *OrderContact) {
	for i, contact := range v.Contacts {
		if contact.ParentID == parentID {
			return i, contact
		}
	}
	return -1, nil
}

// OrderContact is summary of contact information stored in order
// TODO: Try to remove or document why we can't use only OrderCounterparty instead. Probably because of ParentID field?
type OrderContact struct {
	// Required fields
	ID    string                       `json:"id" firestore:"id"`
	Type  briefs4contactus.ContactType `json:"type" firestore:"type"`
	Title string                       `json:"title" firestore:"title"`
	// Optional fields
	ParentID string `json:"parentID,omitempty" firestore:"parentID,omitempty"`

	// CountryID is used instead of address
	// We do not need address in order contact yet as we store it
	// in OrderShippingPoint.Location of ShippingPointLocation.Address type.
	// Address dbmodels.Address `json:"address,omitempty" firestore:"address,omitempty"`
	CountryID string `json:"countryID,omitempty" firestore:"countryID,omitempty"`
}

func (v OrderContact) String() string {
	return fmt.Sprintf(`OrderContact{ID=%s, Type=%s, ParentID=%s, Title=%s}`, v.ID, v.Type, v.ParentID, v.Title)
}

// Validate returns error if OrderContact is not valid
func (v OrderContact) Validate() error {
	if err := briefs4contactus.ValidateContactIDRecordField("id", v.ID, true); err != nil {
		return err
	}
	if err := briefs4contactus.ValidateContactType(v.Type); err != nil {
		return validation.NewErrBadRecordFieldValue("contactID", err.Error())
	}
	if v.ParentID == v.ID {
		return validation.NewErrBadRecordFieldValue("parentID", "parentID cannot be the same as ID")
	}
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	if err := with.ValidateCountryID("countryID", v.CountryID, true); err != nil {
		return err
	}

	//if err := v.Address.Validate(); err != nil {
	//	return validation.NewErrBadRecordFieldValue("address", err.Error())
	//}
	//if v.Type == dbmodels.ContactTypeLocation && len(strings.TrimSpace(v.Address.Lines)) == 0 {
	//	return validation.NewErrRecordIsMissingRequiredField("address.lines")
	//}
	return nil
}
