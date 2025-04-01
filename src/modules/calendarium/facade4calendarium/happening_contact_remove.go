package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func RemoveParticipantFromHappening(ctx facade.ContextWithUser, request dto4calendarium.HappeningContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	if err = dal4calendarium.RunHappeningSpaceWorker(ctx, ctx.User(), request.HappeningRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return removeParticipantFromHappeningTxWorker(ctx, tx, params, request)
	}); err != nil {
		return fmt.Errorf("failed to remove participant from happening: %w", err)
	}
	return nil
}

func removeParticipantFromHappeningTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.HappeningContactRequest) error {
	contactRequest := dto4calendarium.HappeningContactsRequest{
		HappeningRequest: request.HappeningRequest,
		Contacts:         []dbo4linkage.ShortSpaceModuleDocRef{request.Contact},
	}
	_, err := getHappeningContactRecords(ctx, tx, &contactRequest, params)
	if err != nil {
		return err
	}
	contactShortRef := dbmodels.NewSpaceItemID(request.Contact.SpaceID, request.Contact.ID)
	switch params.Happening.Data.Type {
	case dbo4calendarium.HappeningTypeSingle:
		break // nothing to do
	case dbo4calendarium.HappeningTypeRecurring:
		var updates []update.Update
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
	contactFullRef := dbo4contactus.NewContactFullRef(contactShortRef.SpaceID(), contactShortRef.ItemID())
	params.HappeningUpdates = append(
		params.HappeningUpdates,
		dbo4linkage.RemoveRelatedAndID(
			&params.Happening.Data.WithRelated,
			&params.Happening.Data.WithRelatedIDs,
			contactFullRef,
		)...,
	)
	params.Happening.Record.MarkAsChanged()
	return err
}

func removeContactFromHappeningBriefInContactusSpaceDbo(
	calendariumSpace dal4calendarium.CalendariumSpaceEntry,
	happening dbo4calendarium.HappeningEntry,
	contactShortRef dbmodels.SpaceItemID,
) (updates []update.Update, err error) {
	calendarHappeningBrief := calendariumSpace.Data.GetRecurringHappeningBrief(happening.ID)
	if calendarHappeningBrief == nil {
		return nil, err
	}
	contactFullRef := dbo4contactus.NewContactFullRef(contactShortRef.SpaceID(), contactShortRef.ItemID())
	updates = calendarHappeningBrief.RemoveRelatedItem(contactFullRef)
	if len(updates) > 0 {
		for i, u := range updates {
			updates[i] = update.ByFieldName(
				fmt.Sprintf("recurringHappenings.%s.%s", happening.ID, u.FieldName()),
				u.Value(),
			)
		}
	}
	return updates, nil
}
