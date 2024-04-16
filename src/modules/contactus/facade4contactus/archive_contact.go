package facade4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// ArchiveContact archives team contact - e.g., hides it from the list of contacts
func ArchiveContact(ctx context.Context, userContext facade.User, request dto4contactus.ContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	return dal4contactus.RunContactusTeamWorker(ctx, userContext, request.TeamRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusTeamWorkerParams) (err error) {
			return archiveContactTxWorker(ctx, tx, params, request.ContactID)
		},
	)
}

func archiveContactTxWorker(
	ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusTeamWorkerParams,
	contactID string,
) (err error) {
	contact := dal4contactus.NewContactEntry(params.Team.ID, contactID)
	if err = params.GetRecords(ctx, tx, params.UserID, contact.Record); err != nil {
		return err
	}
	contactUpdates := removeContactRoles(params, contact)
	if len(contactUpdates) > 0 {
		if err = tx.Update(ctx, contact.Record.Key(), contactUpdates); err != nil {
			return err
		}
	}
	return err
}
