package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
)

func GetContactByID(ctx context.Context, tx dal.ReadSession, spaceID, contactID string) (contact ContactEntry, err error) {
	contact = NewContactEntry(spaceID, contactID)
	return contact, GetContact(ctx, tx, contact)
}

func GetContact(ctx context.Context, tx dal.ReadSession, contact ContactEntry) error {
	return tx.Get(ctx, contact.Record)
}
