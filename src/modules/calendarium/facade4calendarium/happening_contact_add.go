package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
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

		switch params.Happening.Dbo.Type {
		case models4calendarium.HappeningTypeSingle:
			break // No special processing needed
		case models4calendarium.HappeningTypeRecurring:
			var updates []dal.Update
			if updates, err = addContactToHappeningBriefInTeamDto(ctx, tx, params.TeamModuleEntry, params.Happening, request.Contact.ID); err != nil {
				return fmt.Errorf("failed to add member to happening brief in team DTO: %w", err)
			}
			params.TeamModuleUpdates = append(params.TeamModuleUpdates, updates...)
		default:
			return fmt.Errorf("invalid happenning record: %w",
				validation.NewErrBadRecordFieldValue("type",
					fmt.Sprintf("unknown value: [%v]", params.Happening.Dbo.Type)))
		}
		contactFullRef := models4contactus.NewContactFullRef(request.TeamID, request.Contact.ID)
		var updates []dal.Update
		if updates, err = params.Happening.Dbo.AddRelationshipAndID(
			models4linkage.Link{
				TeamModuleItemRef: contactFullRef,
				RelatedAs:         []string{"participant"},
			},
		); err != nil {
			return err
		}
		params.HappeningUpdates = append(params.HappeningUpdates, updates...)
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
	calendariumTeam dal4calendarium.CalendariumTeamContext,
	happening models4calendarium.HappeningContext,
	contactID string,
) (updates []dal.Update, err error) {
	teamID := calendariumTeam.Key.Parent().ID.(string)
	happeningBriefPointer := calendariumTeam.Data.GetRecurringHappeningBrief(happening.ID)
	//teamContactID := dbmodels.NewTeamItemID(teamID, contactID)
	var happeningBrief models4calendarium.HappeningBrief
	if happeningBriefPointer == nil {
		happeningBrief = happening.Dbo.HappeningBrief // Make copy so we do not affect the DTO object
		happeningBriefPointer = &models4calendarium.CalendarHappeningBrief{
			HappeningBrief: happeningBrief,
			WithRelated:    happening.Dbo.WithRelated,
		}
		//} else if happeningBriefPointer.Participants[string(teamContactID)] != nil {
		//	return nil // Already added to happening brief in calendariumTeam record
	}
	contactRef := models4linkage.NewTeamModuleDocRef(teamID, const4contactus.ModuleID, const4contactus.ContactsCollection, contactID)

	updates, err = happeningBriefPointer.AddRelationship(models4linkage.Link{TeamModuleItemRef: contactRef, RelatedAs: []string{"participant"}})
	for i := range updates {
		updates[i].Field = fmt.Sprintf("recurringHappenings.%s.%s", happening.ID, updates[i].Field)
	}

	//if happeningBriefPointer.Participants == nil {
	//	happeningBriefPointer.Participants = make(map[string]*models4calendarium.HappeningParticipant)
	//}
	//if happeningBriefPointer.Participants[string(teamContactID)] == nil {
	//	happeningBriefPointer.Participants[string(teamContactID)] = &models4calendarium.HappeningParticipant{}
	//}
	//if calendariumTeam.Data.RecurringHappenings == nil {
	//	calendariumTeam.Data.RecurringHappenings = make(map[string]*models4calendarium.CalendarHappeningBrief, 1)
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
