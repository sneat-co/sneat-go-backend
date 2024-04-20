package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func RemoveParticipantFromHappening(ctx context.Context, user facade.User, request dto4calendarium.HappeningContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *happeningWorkerParams) error {
		_, err := getHappeningContactRecords(ctx, tx, &request, params)
		if err != nil {
			return err
		}
		contactShortRef := dbmodels.NewTeamItemID(request.Contact.TeamID, request.Contact.ID)
		switch params.Happening.Dbo.Type {
		case "single":
			break // nothing to do
		case "recurring":
			var updates []dal.Update
			if updates, err = removeContactFromHappeningBriefInContactusTeamDbo(params.TeamModuleEntry, params.Happening, contactShortRef); err != nil {
				return fmt.Errorf("failed to remove member from happening brief in team DBO: %w", err)
			}
			params.TeamModuleUpdates = append(params.TeamModuleUpdates, updates...)
		default:
			return fmt.Errorf("invalid happenning record: %w",
				validation.NewErrBadRecordFieldValue("type",
					fmt.Sprintf("unknown value: [%v]", params.Happening.Dbo.Type)))
		}
		contactFullRef := models4contactus.NewContactFullRef(contactShortRef.TeamID(), contactShortRef.ItemID())
		params.HappeningUpdates = append(params.HappeningUpdates, params.Happening.Dbo.RemoveRelatedAndID(contactFullRef)...)
		return err
	}

	if err = modifyHappening(ctx, user, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}

func removeContactFromHappeningBriefInContactusTeamDbo(
	calendariumTeam dal4calendarium.CalendariumTeamContext,
	happening models4calendarium.HappeningContext,
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
