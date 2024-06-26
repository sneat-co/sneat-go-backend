package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func RemoveParticipantFromHappening(ctx context.Context, user facade.User, request dto4calendarium.HappeningContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return removeParticipantFromHappeningTxWorker(ctx, tx, params, request)
	}

	if err = dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}

func removeParticipantFromHappeningTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.HappeningContactRequest) error {
	_, err := getHappeningContactRecords(ctx, tx, &request, params)
	if err != nil {
		return err
	}
	contactShortRef := dbmodels.NewTeamItemID(request.Contact.TeamID, request.Contact.ID)
	switch params.Happening.Data.Type {
	case dbo4calendarium.HappeningTypeSingle:
		break // nothing to do
	case dbo4calendarium.HappeningTypeRecurring:
		var updates []dal.Update
		if updates, err = removeContactFromHappeningBriefInContactusTeamDbo(params.TeamModuleEntry, params.Happening, contactShortRef); err != nil {
			return fmt.Errorf("failed to remove member from happening brief in team DBO: %w", err)
		}
		params.TeamModuleUpdates = append(params.TeamModuleUpdates, updates...)
		params.TeamModuleEntry.Record.MarkAsChanged()
	default:
		return fmt.Errorf("invalid happenning record: %w",
			validation.NewErrBadRecordFieldValue("type",
				fmt.Sprintf("unknown value: [%v]", params.Happening.Data.Type)))
	}
	contactFullRef := models4contactus.NewContactFullRef(contactShortRef.TeamID(), contactShortRef.ItemID())
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

func removeContactFromHappeningBriefInContactusTeamDbo(
	calendariumTeam dal4calendarium.CalendariumTeamEntry,
	happening dbo4calendarium.HappeningEntry,
	contactShortRef dbmodels.TeamItemID,
) (updates []dal.Update, err error) {
	calendarHappeningBrief := calendariumTeam.Data.GetRecurringHappeningBrief(happening.ID)
	if calendarHappeningBrief == nil {
		return nil, err
	}
	contactFullRef := models4contactus.NewContactFullRef(contactShortRef.TeamID(), contactShortRef.ItemID())
	updates = calendarHappeningBrief.WithRelated.RemoveRelatedItem(contactFullRef)
	if len(updates) > 0 {
		for i := range updates {
			updates[i].Field = fmt.Sprintf("recurringHappenings.%s.%s", happening.ID, updates[i].Field)
		}
	}
	return updates, nil
}
