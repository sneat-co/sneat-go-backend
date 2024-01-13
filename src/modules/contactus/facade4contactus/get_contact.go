package facade4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
)

func GetContactByID(ctx context.Context, tx dal.ReadSession, teamID, contactID string) (contact dal4contactus.ContactEntry, err error) {
	contact = dal4contactus.NewContactEntry(teamID, contactID)
	return contact, tx.Get(ctx, contact.Record)
}
