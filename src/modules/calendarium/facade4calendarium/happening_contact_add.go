package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

func AddParticipantToHappening(ctx context.Context, user facade.User, request dto4calendarium.HappeningContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4calendarium.HappeningWorkerParams) error {
		return addParticipantToHappeningTxWorker(ctx, tx, params, request)
	}

	if err = dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, worker); err != nil {
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
		var updates []dal.Update
		if updates, err = addContactToHappeningBriefInTeamDto(ctx, tx, params.TeamModuleEntry, params.Happening, request.Contact.ID); err != nil {
			return fmt.Errorf("failed to add member to happening brief in team DTO: %w", err)
		}
		params.TeamModuleUpdates = append(params.TeamModuleUpdates, updates...)
		params.TeamModuleEntry.Record.MarkAsChanged()
	default:
		return fmt.Errorf("invalid happenning record: %w",
			validation.NewErrBadRecordFieldValue("type",
				fmt.Sprintf("unknown value: [%v]", params.Happening.Data.Type)))
	}
	contactFullRef := models4contactus.NewContactFullRef(request.TeamID, request.Contact.ID)
	var updates []dal.Update
	if updates, err = dbo4linkage.AddRelationshipAndID(
		&params.Happening.Data.WithRelated,
		&params.Happening.Data.WithRelatedIDs,
		contactFullRef,
		dbo4linkage.RelationshipRolesCommand{
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

func addContactToHappeningBriefInTeamDto(
	_ context.Context,
	_ dal.ReadwriteTransaction,
	calendariumTeam dal4calendarium.CalendariumTeamEntry,
	happening dbo4calendarium.HappeningEntry,
	contactID string,
) (updates []dal.Update, err error) {
	teamID := calendariumTeam.Key.Parent().ID.(string)
	happeningBriefPointer := calendariumTeam.Data.GetRecurringHappeningBrief(happening.ID)
	var happeningBrief dbo4calendarium.HappeningBrief
	if happeningBriefPointer == nil {
		happeningBrief = happening.Data.HappeningBrief // Make copy so we do not affect the DTO object
		happeningBriefPointer = &dbo4calendarium.CalendarHappeningBrief{
			HappeningBrief: happeningBrief,
			WithRelated:    happening.Data.WithRelated,
		}
	}
	contactRef := dbo4linkage.NewTeamModuleItemRef(teamID, const4contactus.ModuleID, const4contactus.ContactsCollection, contactID)

	updates, err = happeningBriefPointer.AddRelationship(
		contactRef,
		dbo4linkage.RelationshipRolesCommand{
			Add: &dbo4linkage.RolesCommand{
				RolesOfItem: []string{"participant"},
			},
		})
	for i := range updates {
		updates[i].Field = fmt.Sprintf("recurringHappenings.%s.%s", happening.ID, updates[i].Field)
	}

	calendariumTeam.Data.RecurringHappenings[happening.ID] = happeningBriefPointer
	return
}
