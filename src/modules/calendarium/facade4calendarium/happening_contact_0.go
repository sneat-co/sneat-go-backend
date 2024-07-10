package facade4calendarium

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
)

func getHappeningContactRecords(ctx context.Context, tx dal.ReadwriteTransaction, request *dto4calendarium.HappeningContactRequest, params *dal4calendarium.HappeningWorkerParams) (contact dal4contactus.ContactEntry, err error) {
	if request.Contact.SpaceID == "" {
		request.Contact.SpaceID = request.SpaceID
	}
	contact = dal4contactus.NewContactEntry(request.Contact.SpaceID, request.Contact.ID)

	if err = tx.GetMulti(ctx, []dal.Record{params.Happening.Record, params.SpaceModuleEntry.Record, contact.Record}); err != nil {
		return contact, fmt.Errorf("failed to get records: %w", err)
	}
	if err = params.SpaceModuleEntry.Record.Error(); err != nil {
		if !dal.IsNotFound(err) && !errors.Is(err, dal.NoError) {
			return contact, fmt.Errorf("failed to get contactus team record: %w", err)
		}
	}
	if !params.SpaceModuleEntry.Record.Exists() {
		return contact, fmt.Errorf("happening not found: %w", params.SpaceModuleEntry.Record.Error())
	}
	if !contact.Record.Exists() {
		return contact, fmt.Errorf("contact not found: %w", contact.Record.Error())
	}
	return
}
