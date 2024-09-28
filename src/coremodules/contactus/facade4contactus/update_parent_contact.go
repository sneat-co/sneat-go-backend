package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/strongo/logus"
)

func updateParentContact(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	contact, parent dal4contactus.ContactEntry,
) error {
	//logus.Debugf(ctx, "updateParentContact(contact=%s, parentID=%s)", contact.ContactID, parent.ContactID)
	contactBrief := &briefs4contactus.ContactBrief{
		Type:   contact.Data.Type,
		Gender: contact.Data.Gender,
		Names:  contact.Data.Names,
	}
	contactBrief.RelatedAs = RelatedAsChild
	updates := parent.Data.SetContactBrief(contact.Key.Parent().ID.(string), contact.ID, contactBrief)
	if err := parent.Data.Validate(); err != nil {
		return fmt.Errorf("parent contact DBO validation error: %w", err)
	}
	if err := tx.Update(ctx, parent.Key, updates); err != nil {
		return fmt.Errorf("failed to update parent contact record: %w", err)
	}
	logus.Infof(ctx, "updateParentContact(contact=%v, parentID=%v) - success!", contact.ID, parent.ID)
	return nil
}
