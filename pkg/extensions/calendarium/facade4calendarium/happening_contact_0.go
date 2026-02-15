package facade4calendarium

import (
	"context"
	"errors"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dto4calendarium"
)

func getHappeningContactRecords(ctx context.Context, tx dal.ReadwriteTransaction, request *dto4calendarium.HappeningContactsRequest, params *dal4calendarium.HappeningWorkerParams) (contacts []dal4contactus.ContactEntry, err error) {
	records := make([]dal.Record, 0, len(request.Contacts)+2)
	records = append(records, params.Happening.Record)
	records = append(records, params.SpaceModuleEntry.Record)

	for _, contactRef := range request.Contacts {
		if contactRef.SpaceID == "" {
			contactRef.SpaceID = request.SpaceID
		}
		contact := dal4contactus.NewContactEntry(contactRef.SpaceID, contactRef.ID)
		contacts = append(contacts, contact)
		records = append(records, contact.Record)
	}
	if err = tx.GetMulti(ctx, records); err != nil {
		return contacts, fmt.Errorf("failed to get records: %w", err)
	}
	if err = params.SpaceModuleEntry.Record.Error(); err != nil {
		if !dal.IsNotFound(err) && !errors.Is(err, dal.ErrNoError) {
			return contacts, fmt.Errorf("failed to get contactus team record: %w", err)
		}
	}
	if !params.SpaceModuleEntry.Record.Exists() {
		return contacts, fmt.Errorf("happening not found: %w", params.SpaceModuleEntry.Record.Error())
	}
	for _, contact := range contacts {
		if !contact.Record.Exists() {
			return contacts, fmt.Errorf("contact not found: %w", contact.Record.Error())
		}
	}
	return
}
