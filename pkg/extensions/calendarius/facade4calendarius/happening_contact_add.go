package facade4calendarius

import (
	"context"
	"fmt"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dal4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

func AddParticipantsToHappening(ctx facade.ContextWithUser, request dto4calendarius.HappeningContactsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	if err = dal4calendarius.RunHappeningSpaceWorker(ctx, request.HappeningRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams) error {
			return addParticipantToHappeningTxWorker(ctx, tx, params, request)
		}); err != nil {
		return fmt.Errorf("failed to add participant to happening: %w", err)
	}
	return nil
}

func addParticipantToHappeningTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarius.HappeningWorkerParams, request dto4calendarius.HappeningContactsRequest) error {
	_, err := getHappeningContactRecords(ctx, tx, &request, params)
	if err != nil {
		return err
	}

	switch params.Happening.Data.Type {
	case dbo4calendarius.HappeningTypeSingle:
		break // No special processing needed
	case dbo4calendarius.HappeningTypeRecurring:
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
	for _, contactShortRef := range request.Contacts {
		if contactShortRef.SpaceID == request.SpaceID {
			contactShortRef.SpaceID = ""
		}
		var contactRef dbo4linkage.ItemRef
		if contactShortRef.SpaceID == "" {
			contactRef = dbo4contactus.NewContactRefSameSpace(contactShortRef.ID)
		} else {
			contactRef = dbo4contactus.NewContactFullRef(contactShortRef.SpaceID, contactShortRef.ID)
		}
		if updates, err = dbo4linkage.AddRelationshipAndID(
			params.Started,
			params.UserID(),
			params.Space.ID,
			&params.Happening.Data.WithRelated,
			&params.Happening.Data.WithRelatedIDs,
			dbo4linkage.RelationshipItemRolesCommand{
				ItemRef: contactRef,
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
	calendariusSpace dal4calendarius.CalendariusSpaceEntry,
	happening dbo4calendarius.HappeningEntry,
	contactRefs []dbo4linkage.ShortSpaceModuleItemRef,
) (updates []update.Update, err error) {
	if len(contactRefs) == 0 {
		return updates, fmt.Errorf("no contacts to add to happening")
	}
	spaceID := coretypes.SpaceID(calendariusSpace.Key.Parent().ID.(string))
	happeningBriefPointer := calendariusSpace.Data.GetRecurringHappeningBrief(happening.ID)
	var happeningBase dbo4calendarius.HappeningBase
	if happeningBriefPointer == nil {
		happeningBase = happening.Data.HappeningBase // Make copy so we do not affect the DTO object
		happeningBriefPointer = &dbo4calendarius.CalendarHappeningBrief{
			HappeningBase: happeningBase,
			WithRelated:   happening.Data.WithRelated,
		}
	}
	for _, contactShortRef := range contactRefs {
		if contactShortRef.SpaceID == spaceID {
			contactShortRef.SpaceID = ""
		}
		var contactRef dbo4linkage.ItemRef
		if contactShortRef.SpaceID == "" {
			contactRef = dbo4contactus.NewContactRefSameSpace(contactShortRef.ID)
		} else {
			contactRef = dbo4contactus.NewContactFullRef(contactShortRef.SpaceID, contactShortRef.ID)
		}

		updates, err = happeningBriefPointer.ProcessRelatedCommand(
			now,
			userID,
			dbo4linkage.RelationshipItemRolesCommand{
				ItemRef: contactRef,
				Add: &dbo4linkage.RolesCommand{
					RolesOfItem: []string{"participant"},
				},
			})
	}

	hFieldPath := []string{dbo4calendarius.RecurringHappeningsField, happening.ID}
	for i, u := range updates {
		v := u.Value()
		if fieldPath := u.FieldPath(); len(fieldPath) == 0 {
			updates[i] = update.ByFieldPath(append(hFieldPath, u.FieldName()), v)
		} else {
			updates[i] = update.ByFieldPath(append(hFieldPath, fieldPath...), v)
		}
	}
	calendariusSpace.Data.RecurringHappenings[happening.ID] = happeningBriefPointer
	return
}
