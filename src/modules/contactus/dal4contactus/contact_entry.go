package dal4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
)

type ContactEntry = record.DataWithID[string, *models4contactus.ContactDbo]

func NewContactEntry(teamID, contactID string) ContactEntry {
	return NewContactEntryWithData(teamID, contactID, new(models4contactus.ContactDbo))
}

func NewContactEntryWithData(teamID, contactID string, data *models4contactus.ContactDbo) (contact ContactEntry) {
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
