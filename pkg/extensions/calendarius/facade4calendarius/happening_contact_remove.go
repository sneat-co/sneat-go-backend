package facade4calendarius

import (
	"context"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dal4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func RemoveParticipantsFromHappening(ctx facade.ContextWithUser, request dto4calendarius.HappeningContactsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	if err = dal4calendarius.RunHappeningSpaceWorker(ctx, request.HappeningRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams) error {
			return removeParticipantsFromHappeningTxWorker(ctx, tx, params, request)
		}); err != nil {
		return fmt.Errorf("failed to remove participant from happening: %w", err)
	}
	return nil
}

func removeParticipantsFromHappeningTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams, request dto4calendarius.HappeningContactsRequest) error {
	for i := range request.Contacts {
		// Contacts splice is holding non-pointer structs so we need to use index
		if request.Contacts[i].SpaceID == "" {
			request.Contacts[i].SpaceID = request.SpaceID
		}
	}
	_, err := getHappeningContactRecords(ctx, tx, &request, params)
	if err != nil {
		return err
	}
	for _, contact := range request.Contacts {
		switch params.Happening.Data.Type {
		case dbo4calendarius.HappeningTypeSingle:
			// nothing to do
		case dbo4calendarius.HappeningTypeRecurring:
			var updates []update.Update

			contactShortRef := dbmodels.NewSpaceItemID(contact.SpaceID, contact.ID)
			if updates, err = removeContactFromHappeningBriefInContactusSpaceDbo(params.SpaceModuleEntry, params.Happening, contactShortRef); err != nil {
				return fmt.Errorf("failed to remove member from happening brief in space DBO: %w", err)
			}
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, updates...)
			params.SpaceModuleEntry.Record.MarkAsChanged()
		default:
			return fmt.Errorf("invalid happenning record: %w",
				validation.NewErrBadRecordFieldValue("type",
					fmt.Sprintf("unknown value: [%v]", params.Happening.Data.Type)))
		}

		var contactRef dbo4linkage.ItemRef
		if contact.SpaceID == "" || contact.SpaceID == request.SpaceID {
			contactRef = dbo4contactus.NewContactRefSameSpace(contact.ID)
		} else {
			contactRef = dbo4contactus.NewContactFullRef(contact.SpaceID, contact.ID)
		}
		params.HappeningUpdates = append(
			params.HappeningUpdates,
			dbo4linkage.RemoveRelatedAndID(
				params.Space.ID,
				&params.Happening.Data.WithRelated,
				&params.Happening.Data.WithRelatedIDs,
				contactRef,
			)...,
		)
		params.Happening.Record.MarkAsChanged()
	}

	return err
}

func removeContactFromHappeningBriefInContactusSpaceDbo(
	calendariusSpace dal4calendarius.CalendariusSpaceEntry,
	happening dbo4calendarius.HappeningEntry,
	contactShortRef dbmodels.SpaceItemID,
) (updates []update.Update, err error) {
	calendarHappeningBrief := calendariusSpace.Data.GetRecurringHappeningBrief(happening.ID)
	if calendarHappeningBrief == nil {
		return nil, err
	}
	var contactRef dbo4linkage.ItemRef
	if contactSpaceID := contactShortRef.SpaceID(); contactSpaceID == "" || contactSpaceID == coretypes.SpaceID(calendariusSpace.Key.Parent().ID.(string)) {
		contactRef = dbo4contactus.NewContactRefSameSpace(contactShortRef.ItemID())
	} else {
		contactRef = dbo4contactus.NewContactFullRef(contactShortRef.SpaceID(), contactShortRef.ItemID())
	}
	updates = calendarHappeningBrief.RemoveRelatedItem(contactRef)
	if len(updates) > 0 {
		for i, u := range updates {
			fieldPath := []string{dbo4calendarius.RecurringHappeningsField, happening.ID}
			if fieldName := u.FieldName(); fieldName != "" {
				fieldPath = append(fieldPath, fieldName)
			} else {
				fieldPath = append(fieldPath, u.FieldPath()...)
			}
			updates[i] = update.ByFieldPath(fieldPath, u.Value())
		}
	}
	return updates, nil
}
