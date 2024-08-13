package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
)

type ContactEntry = record.DataWithID[string, *models4contactus.ContactDbo]

func NewContactEntry(spaceID, contactID string) ContactEntry {
	return NewContactEntryWithData(spaceID, contactID, new(models4contactus.ContactDbo))
}

func NewContactEntryWithData(spaceID, contactID string, data *models4contactus.ContactDbo) (contact ContactEntry) {
	key := NewContactKey(spaceID, contactID)
	contact.ID = contactID
	contact.FullID = spaceID + ":" + contactID
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

func GetContactusSpace(ctx context.Context, tx dal.ReadSession, contactusSpace ContactusSpaceEntry) (err error) {
	return tx.Get(ctx, contactusSpace.Record)
}
