package briefs4contactus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"slices"
)

type contactBrief interface {
	dbmodels.UserIDGetter
	dbmodels.RelatedAs
}

type WithSingleSpaceContactsWithoutContactIDs[
	T interface {
		contactBrief
		HasRole(role string) bool
		Equal(v T) bool
	},
] struct {
	WithContactsBase[T]
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) Validate() error {
	for id, brief := range v.Contacts {
		if err := validate.RecordID(id); err != nil {
			return validation.NewErrBadRecordFieldValue(const4contactus.ContactsField,
				fmt.Sprintf("invalid contact ContactID=%s: %v", id, err))
		}
		if err := brief.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("contacts."+id, err.Error())
		}
	}
	return nil
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) ContactIDs() (contactIDs []string) {
	contactIDs = make([]string, 0, len(v.Contacts))
	for id := range v.Contacts {
		contactIDs = append(contactIDs, id)
	}
	return
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) HasContact(contactID string) bool {
	_, ok := v.Contacts[contactID]
	return ok
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) AddContact(contactID string, contact T) dal.Update {
	if v.Contacts == nil {
		v.Contacts = make(map[string]T)
	}
	v.Contacts[contactID] = contact
	return dal.Update{
		Field: "contacts." + contactID,
		Value: contact,
	}
}

func (v *WithSingleSpaceContactsWithoutContactIDs[T]) RemoveContact(contactID string) dal.Update {
	delete(v.Contacts, contactID)
	return dal.Update{
		Field: "contacts." + contactID,
		Value: dal.DeleteField,
	}
}

// WithMultiSpaceContacts mixin that adds WithMultiSpaceContactIDs.ContactIDs & Contacts fields
type WithMultiSpaceContacts[
	T interface {
		contactBrief
		HasRole(role string) bool
		Equal(v T) bool
	},
] struct {
	WithMultiSpaceContactIDs
	WithContactsBase[T]
}

// Validate returns error if not valid
func (v *WithMultiSpaceContacts[T]) Validate() error {
	if err := v.WithMultiSpaceContactIDs.Validate(); err != nil {
		return nil
	}
	return dbmodels.ValidateWithIdsAndBriefs("contactIDs", const4contactus.ContactsField, v.ContactIDs, v.Contacts)
}

func (v *WithMultiSpaceContacts[T]) Updates(contactIDs ...dbmodels.SpaceItemID) (updates []dal.Update) {
	updates = append(updates,
		dal.Update{
			Field: "contactIDs",
			Value: v.ContactIDs,
		},
	)
	if len(contactIDs) == 0 {
		updates = append(updates, dal.Update{
			Field: const4contactus.ContactsField,
			Value: v.Contacts,
		})
	} else {
		for _, id := range contactIDs {
			updates = append(updates, dal.Update{
				Field: const4contactus.ContactsField + "." + string(id),
				Value: v.Contacts[string(id)],
			})
		}
	}
	return
}

// SetContactBrief sets contactBrief brief by ContactID
func (v *WithMultiSpaceContacts[T]) SetContactBrief(teamID, contactID string, contactBrief T) (updates []dal.Update) {
	id := string(dbmodels.NewSpaceItemID(teamID, contactID))
	if !slices.Contains(v.ContactIDs, id) {
		v.ContactIDs = append(v.ContactIDs, id)
		updates = append(updates, dal.Update{
			Field: "contactIDs",
			Value: v.ContactIDs,
		})
	}
	if currentBrief, ok := v.Contacts[id]; !ok || !currentBrief.Equal(contactBrief) {
		v.Contacts[id] = contactBrief
		updates = append(updates, dal.Update{
			Field: const4contactus.ContactsField + "." + string(id),
			Value: contactBrief,
		})
	}
	return
}

// ParentContactBrief returns parent contactBrief brief
func (v *WithMultiSpaceContacts[T]) ParentContactBrief() (i int, id dbmodels.SpaceItemID, brief T) {
	for i, id := range v.ContactIDs {
		brief := v.Contacts[id]
		if brief.GetRelatedAs() == "parent" {
			return i, dbmodels.SpaceItemID(id), brief
		}
	}
	return -1, "", brief
}

// GetContactBriefByID returns contactBrief brief by ContactID
func (v *WithMultiSpaceContacts[T]) GetContactBriefByID(teamID, contactID string) (i int, brief T) {
	id := dbmodels.NewSpaceItemID(teamID, contactID)
	if brief, ok := v.Contacts[string(id)]; !ok {
		return -1, brief
	}
	return slice.Index(v.ContactIDs, string(id)), brief
}

// GetContactBriefByUserID returns contactBrief brief by user ContactID
func (v *WithMultiSpaceContacts[T]) GetContactBriefByUserID(userID string) (id dbmodels.SpaceItemID, t T) {
	for cID, c := range v.Contacts {
		if c.GetUserID() == userID {
			return dbmodels.SpaceItemID(cID), c
		}
	}
	return
}

func (v *WithMultiSpaceContacts[T]) AddContact(teamID, contactID string, c T) (updates []dal.Update) {
	id := dbmodels.NewSpaceItemID(teamID, contactID)
	if !slices.Contains(v.ContactIDs, string(id)) {
		if len(v.ContactIDs) == 0 {
			v.ContactIDs = make([]string, 1, 2)
			v.ContactIDs[0] = "*"
		}
		v.ContactIDs = append(v.ContactIDs, string(id))
		updates = append(updates, dal.Update{
			Field: "contactIDs",
			Value: v.ContactIDs,
		})
	}
	if _, ok := v.Contacts[string(id)]; !ok {
		updates = append(updates, dal.Update{
			Field: "contacts." + string(id),
			Value: c,
		})
	}
	if v.Contacts == nil {
		v.Contacts = make(map[string]T)
	}
	v.Contacts[string(id)] = c
	return
}

func (v *WithMultiSpaceContacts[T]) RemoveContact(teamID, contactID string) (updates []dal.Update) {
	id := dbmodels.NewSpaceItemID(teamID, contactID)
	contactIDs := slice.RemoveInPlace(string(id), v.ContactIDs)
	if len(contactIDs) != len(v.ContactIDs) {
		v.ContactIDs = contactIDs
		updates = append(updates, dal.Update{
			Field: "contactIDs",
			Value: v.ContactIDs,
		})
	}
	if _, ok := v.Contacts[string(id)]; ok {
		delete(v.Contacts, string(id))
		updates = append(updates, dal.Update{
			Field: "contacts." + string(id),
			Value: dal.DeleteField,
		})
	}
	return
}
