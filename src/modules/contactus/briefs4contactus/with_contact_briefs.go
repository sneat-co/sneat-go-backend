package briefs4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// WithContactBriefs is a base struct for DTOs that have contacts
// Unlike the WithContactsBase it does not keep UserIDs []string.
type WithContactBriefs[
	T interface {
		core.Validatable
		Equal(v T) bool
	},
] struct {
	Contacts map[string]T `json:"contacts,omitempty" firestore:"contacts,omitempty"`
}

type ContactBriefer interface {
	core.Validatable
	dbmodels.UserIDGetter
	dbmodels.RelatedAs
	HasRole(role string) bool
}

// WithContactsBase is a base struct for DTOs that represent a short version of a contact
// TODO: Document how it is different from WithContactBriefs or merge them
type WithContactsBase[T interface {
	ContactBriefer
	Equal(v T) bool
}] struct {
	WithContactBriefs[T]
	dbmodels.WithUserIDs
}

func (v WithContactsBase[T]) Validate() error {
	for id, contact := range v.Contacts {
		if err := contact.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("contacts."+string(id), err.Error())
		}
		if userID := contact.GetUserID(); userID == "" {
			if !v.HasUserID(userID) {
				return validation.NewErrBadRecordFieldValue(
					fmt.Sprintf("contacts.%s.userID", id),
					fmt.Sprintf("%s not added to userIDs", userID))
			}
		}
	}
	return nil
}

func (v WithContactsBase[T]) GetContactBriefByUserID(userID string) (id string, contactBrief T) {
	for k, c := range v.Contacts {
		if c.GetUserID() == userID {
			return k, c
		}
	}
	return id, contactBrief
}

func (v WithContactsBase[T]) GetContactBriefByContactID(contactID string) (contactBrief T) {
	return v.Contacts[contactID]
}

func (v WithContactsBase[T]) GetContactBriefsByRoles(roles ...string) map[string]T {
	result := make(map[string]T)
	for id, c := range v.Contacts {
		for _, role := range roles {
			if c.HasRole(role) {
				result[id] = c
				break
			}
		}
	}
	return result
}

func (v WithContactsBase[T]) GetContactsCount(roles ...string) (count int) {
	for _, c := range v.Contacts {
		for _, role := range roles {
			if c.HasRole(role) {
				count++
				break
			}
		}
	}
	return count
}
