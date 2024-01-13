package dal4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
)

type ContactEntry = record.DataWithID[string, *models4contactus.ContactDto]

func NewContactEntry(teamID, contactID string) ContactEntry {
	return NewContactEntryWithData(teamID, contactID, new(models4contactus.ContactDto))
}

func NewContactEntryWithData(teamID, contactID string, data *models4contactus.ContactDto) (contact ContactEntry) {
	key := NewContactKey(teamID, contactID)
	contact.ID = contactID
	contact.FullID = teamID + ":" + contactID
	contact.Key = key
	contact.Data = data
	contact.Record = dal.NewRecordWithData(key, data)
	return
}

func FindContactEntryByContactID(contacts []ContactEntry, contactID string) (contact ContactEntry, found bool) {
	for _, contact := range contacts {
		if contact.ID == contactID {
			return contact, true
		}
	}
	return contact, false
}

func ContactIDs(contacts []ContactEntry) []string {
	ids := make([]string, len(contacts))
	for i, contact := range contacts {
		ids[i] = contact.ID
	}
	return ids
}
