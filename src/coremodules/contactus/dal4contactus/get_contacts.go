package dal4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
)

func NewContacts(spaceID string, contactIDs ...string) (contacts []ContactEntry) {
	contacts = make([]ContactEntry, len(contactIDs))
	for i, id := range contactIDs {
		if id == "" {
			panic(fmt.Sprintf("contactIDs[%d] == 0", i))
		}
		contacts[i] = NewContactEntry(spaceID, id)
	}
	return
}

func ContactRecords(contacts []ContactEntry) (records []dal.Record) {
	records = make([]dal.Record, len(contacts))
	for i, contact := range contacts {
		records[i] = contact.Record
	}
	return
}

func GetContactsByIDs(ctx context.Context, tx dal.ReadSession, spaceID string, contactsIDs []string) (contacts []ContactEntry, err error) {
	if tx == nil {
		if tx, err = facade.GetSneatDB(ctx); err != nil {
			return
		}
	}
	contacts = NewContacts(spaceID, contactsIDs...)
	records := ContactRecords(contacts)
	return contacts, tx.GetMulti(ctx, records)
}
