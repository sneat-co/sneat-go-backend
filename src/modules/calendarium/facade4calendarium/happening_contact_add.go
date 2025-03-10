package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

func AddParticipantToHappening(ctx facade.ContextWithUser, request dto4calendarium.HappeningContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	if err = dal4calendarium.RunHappeningSpaceWorker(ctx, ctx.User(), request.HappeningRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return addParticipantToHappeningTxWorker(ctx, tx, params, request)
	}); err != nil {
		return fmt.Errorf("failed to add participant to happening: %w", err)
	}
	return nil
}

func addParticipantToHappeningTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams, request dto4calendarium.HappeningContactRequest) error {
	_, err := getHappeningContactRecords(ctx, tx, &request, params)
	if err != nil {
		return err
	}

	switch params.Happening.Data.Type {
	case dbo4calendarium.HappeningTypeSingle:
		break // No special processing needed
	case dbo4calendarium.HappeningTypeRecurring:
		var updates []update.Update
		if updates, err = addContactToHappeningBriefInSpaceDbo(ctx, tx, params.SpaceModuleEntry, params.Happening, request.Contact.ID); err != nil {
			return fmt.Errorf("failed to add member to happening brief in team DTO: %w", err)
		}
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, updates...)
		params.SpaceModuleEntry.Record.MarkAsChanged()
	default:
		return fmt.Errorf("invalid happenning record: %w",
			validation.NewErrBadRecordFieldValue("type",
				fmt.Sprintf("unknown value: [%v]", params.Happening.Data.Type)))
	}
	contactFullRef := dbo4contactus.NewContactFullRef(request.SpaceID, request.Contact.ID)
	var updates []update.Update
	if updates, err = dbo4linkage.AddRelationshipAndID(
		&params.Happening.Data.WithRelated,
		&params.Happening.Data.WithRelatedIDs,
		contactFullRef,
		dbo4linkage.RelationshipItemRolesCommand{
			Add: &dbo4linkage.RolesCommand{
				RolesOfItem: []string{"participant"},
			},
		},
	); err != nil {
		return err
	}
	params.HappeningUpdates = append(params.HappeningUpdates, updates...)
	params.Happening.Record.MarkAsChanged()

	return err
}

func addContactToHappeningBriefInSpaceDbo(
	_ context.Context,
	_ dal.ReadwriteTransaction,
	calendariumSpace dal4calendarium.CalendariumSpaceEntry,
	happening dbo4calendarium.HappeningEntry,
	contactID string,
) (updates []update.Update, err error) {
	spaceID := calendariumSpace.Key.Parent().ID.(coretypes.SpaceID)
	happeningBriefPointer := calendariumSpace.Data.GetRecurringHappeningBrief(happening.ID)
	var happeningBrief dbo4calendarium.HappeningBrief
	if happeningBriefPointer == nil {
		happeningBrief = happening.Data.HappeningBrief // Make copy so we do not affect the DTO object
		happeningBriefPointer = &dbo4calendarium.CalendarHappeningBrief{
			HappeningBrief: happeningBrief,
			WithRelated:    happening.Data.WithRelated,
		}
	}
	contactRef := dbo4linkage.NewSpaceModuleItemRef(spaceID, const4contactus.ModuleID, const4contactus.ContactsCollection, contactID)

	updates, err = happeningBriefPointer.AddRelationship(
		contactRef,
		dbo4linkage.RelationshipItemRolesCommand{
			Add: &dbo4linkage.RolesCommand{
				RolesOfItem: []string{"participant"},
			},
		})
	for i, u := range updates {
		updates[i] = update.ByFieldName(
			fmt.Sprintf("recurringHappenings.%s.%s", happening.ID, u.FieldName()),
			u.Value(),
		)
	}

	calendariumSpace.Data.RecurringHappenings[happening.ID] = happeningBriefPointer
	return
}
