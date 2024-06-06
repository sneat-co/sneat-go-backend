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
		default:
			return fmt.Errorf("invalid happenning record: %w",
				validation.NewErrBadRecordFieldValue("type",
					fmt.Sprintf("unknown value: [%v]", params.Happening.Data.Type)))
		}
		contactFullRef := models4contactus.NewContactFullRef(request.TeamID, request.Contact.ID)
		var updates []dal.Update
		if updates, err = params.Happening.Data.AddRelationshipAndID(
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

		//if params.Happening.Data.ExtraType == dbo4calendarium.HappeningTypeRecurring {
		//	recurringHappening := params.TeamModuleEntry.Data.RecurringHappenings[params.Happening.ID]
		//	if recurringHappening != nil {
		//		recurringHappening.Related = params.Happening.Data.Related
		//		if err = recurringHappening.Validate(); err != nil {
		//			return fmt.Errorf("failed to validate recurring happening: %w", err)
		//		}
		//		if err = params.TeamModuleEntry.Data.Validate(); err != nil {
		//			return fmt.Errorf("failed to validate calendarium team module data: %w", err)
		//		}
		//		params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{
		//			Field: fmt.Sprintf("recurringHappenings.%s.related", params.Happening.ID),
		//		})
		//	}
		//}

		return err
	}

	if err = dal4calendarium.RunHappeningTeamWorker(ctx, user, request.HappeningRequest, worker); err != nil {
		return fmt.Errorf("failed to add member to happening: %w", err)
	}
	return nil
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
	//teamContactID := dbmodels.NewTeamItemID(teamID, contactID)
	var happeningBrief dbo4calendarium.HappeningBrief
	if happeningBriefPointer == nil {
		happeningBrief = happening.Data.HappeningBrief // Make copy so we do not affect the DTO object
		happeningBriefPointer = &dbo4calendarium.CalendarHappeningBrief{
			HappeningBrief: happeningBrief,
			WithRelated:    happening.Data.WithRelated,
		}
		//} else if happeningBriefPointer.Participants[string(teamContactID)] != nil {
		//	return nil // Already added to happening brief in calendariumTeam record
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

	//if happeningBriefPointer.Participants == nil {
	//	happeningBriefPointer.Participants = make(map[string]*dbo4calendarium.HappeningParticipant)
	//}
	//if happeningBriefPointer.Participants[string(teamContactID)] == nil {
	//	happeningBriefPointer.Participants[string(teamContactID)] = &dbo4calendarium.HappeningParticipant{}
	//}
	//if calendariumTeam.Data.RecurringHappenings == nil {
	//	calendariumTeam.Data.RecurringHappenings = make(map[string]*dbo4calendarium.CalendarHappeningBrief, 1)
	//}
	calendariumTeam.Data.RecurringHappenings[happening.ID] = happeningBriefPointer
	//teamUpdates := []dal.Update{
	//	{
	//		Field: "recurringHappenings." + happening.ID,
	//		Value: happeningBriefPointer,
	//	},
	//}
	//if err = tx.Update(ctx, calendariumTeam.Key, teamUpdates); err != nil {
	//	return fmt.Errorf("failed to update calendariumTeam record with a member added to a recurring happening: %w", err)
	//}
	return
}
