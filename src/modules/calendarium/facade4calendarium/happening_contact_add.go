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
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"time"
)

func AddParticipantsToHappening(ctx facade.ContextWithUser, request dto4calendarium.HappeningContactsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	if err = dal4calendarium.RunHappeningSpaceWorker(ctx, request.HappeningRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
			return addParticipantToHappeningTxWorker(ctx, tx, params, request)
		}); err != nil {
		return fmt.Errorf("failed to add participant to happening: %w", err)
	}
	return nil
}

func addParticipantToHappeningTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.HappeningContactsRequest) error {
	_, err := getHappeningContactRecords(ctx, tx, &request, params)
	if err != nil {
		return err
	}

	switch params.Happening.Data.Type {
	case dbo4calendarium.HappeningTypeSingle:
		break // No special processing needed
	case dbo4calendarium.HappeningTypeRecurring:
		var updates []update.Update
		if updates, err = addContactsToHappeningBriefInSpaceDbo(ctx, tx, params.Started, params.UserID(), params.SpaceModuleEntry, params.Happening, request.Contacts); err != nil {
			return fmt.Errorf("failed to add member to happening brief in team DTO: %w", err)
		}
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, updates...)
		params.SpaceModuleEntry.Record.MarkAsChanged()
	default:
		return fmt.Errorf("invalid happenning record: %w",
			validation.NewErrBadRecordFieldValue("type",
				fmt.Sprintf("unknown value: [%v]", params.Happening.Data.Type)))
	}

	var updates []update.Update
	for _, contactRef := range request.Contacts {
		if contactRef.SpaceID == "" {
			contactRef.SpaceID = request.SpaceID
		}
		contactFullRef := dbo4contactus.NewContactFullRef(contactRef.SpaceID, contactRef.ID)
		if updates, err = dbo4linkage.AddRelationshipAndID(
			params.Started,
			params.UserID(),
			params.Space.ID,
			&params.Happening.Data.WithRelated,
			&params.Happening.Data.WithRelatedIDs,
			dbo4linkage.RelationshipItemRolesCommand{
				ItemRef: contactFullRef,
				Add: &dbo4linkage.RolesCommand{
					RolesOfItem: []string{"participant"},
				},
			},
		); err != nil {
			return err
		}
	}
	if len(updates) > 0 {
		params.HappeningUpdates = append(params.HappeningUpdates, updates...)
		params.Happening.Record.MarkAsChanged()
	}

	return err
}

func addContactsToHappeningBriefInSpaceDbo(
	_ context.Context,
	_ dal.ReadwriteTransaction,
	now time.Time,
	userID string,
	calendariumSpace dal4calendarium.CalendariumSpaceEntry,
	happening dbo4calendarium.HappeningEntry,
	contactRefs []dbo4linkage.ShortSpaceModuleItemRef,
) (updates []update.Update, err error) {
	if len(contactRefs) == 0 {
		return updates, fmt.Errorf("no contacts to add to happening")
	}
	spaceID := coretypes.SpaceID(calendariumSpace.Key.Parent().ID.(string))
	happeningBriefPointer := calendariumSpace.Data.GetRecurringHappeningBrief(happening.ID)
	var happeningBase dbo4calendarium.HappeningBase
	if happeningBriefPointer == nil {
		happeningBase = happening.Data.HappeningBase // Make copy so we do not affect the DTO object
		happeningBriefPointer = &dbo4calendarium.CalendarHappeningBrief{
			HappeningBase: happeningBase,
			WithRelated:   happening.Data.WithRelated,
		}
	}
	for _, contactRef := range contactRefs {
		fullContactRef := dbo4contactus.NewContactFullRef(spaceID, contactRef.ID)

		updates, err = happeningBriefPointer.ProcessRelatedCommand(
			now,
			userID,
			dbo4linkage.RelationshipItemRolesCommand{
				ItemRef: fullContactRef,
				Add: &dbo4linkage.RolesCommand{
					RolesOfItem: []string{"participant"},
				},
			})
	}

	hFieldPath := []string{"recurringHappenings", happening.ID}
	for i, u := range updates {
		v := u.Value()
		if fieldPath := u.FieldPath(); len(fieldPath) == 0 {
			updates[i] = update.ByFieldPath(append(hFieldPath, u.FieldName()), v)
		} else {
			updates[i] = update.ByFieldPath(append(hFieldPath, fieldPath...), v)
		}
	}
	calendariumSpace.Data.RecurringHappenings[happening.ID] = happeningBriefPointer
	return
}
