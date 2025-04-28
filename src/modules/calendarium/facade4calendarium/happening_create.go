package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	dal4contactus2 "github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	dal4spaceus2 "github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"strings"
)

// CreateHappening creates a recurring happening
func CreateHappening(
	ctx facade.ContextWithUser, request dto4calendarium.CreateHappeningRequest,
) (
	response dto4calendarium.CreateHappeningResponse, err error,
) {
	request.Happening.Title = strings.TrimSpace(request.Happening.Title)
	if err = request.Validate(); err != nil {
		return
	}
	//var counter string
	//if request.Happening.ExtraType == dbo4calendarium.HappeningTypeRecurring {
	//	counter = "recurringHappenings"
	//}
	happeningDto := &dbo4calendarium.HappeningDbo{
		HappeningBrief: *request.Happening,
		CreatedFields: with.CreatedFields{
			CreatedByField: with.CreatedByField{
				CreatedBy: ctx.User().GetUserID(),
			},
		},
		//WithTeamDates: dbmodels.WithSpaceDates{
		//	WithTeamIDs: dbmodels.WithSpaceIDs{
		//		SpaceIDs: []string{request.Space},
		//	},
		//},
	}
	//happeningDto.ContactIDs = append(happeningDto.ContactIDs, "*")

	if happeningDto.Type == dbo4calendarium.HappeningTypeSingle {
		for _, slot := range happeningDto.Slots {
			date := slot.Start.Date
			if slice.Index(happeningDto.Dates, date) < 0 {
				happeningDto.Dates = append(happeningDto.Dates, date)
			}
		}
	}
	err = dal4spaceus2.CreateSpaceItem(ctx, request.SpaceRequest,
		const4calendarium.ModuleID,
		new(dbo4calendarium.CalendariumSpaceDbo),
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4spaceus2.ModuleSpaceWorkerParams[*dbo4calendarium.CalendariumSpaceDbo]) (err error) {
			response, err = createHappeningTx(ctx, tx, happeningDto, params)
			return
		},
	)
	response.Dbo = *happeningDto
	return
}

func createHappeningTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	happeningDto *dbo4calendarium.HappeningDbo,
	params *dal4spaceus2.ModuleSpaceWorkerParams[*dbo4calendarium.CalendariumSpaceDbo],
) (
	response dto4calendarium.CreateHappeningResponse, err error,
) {
	happeningDto.CreatedAt = params.Started
	contactusSpace := dal4contactus2.NewContactusSpaceEntry(params.Space.ID)
	if err = params.GetRecords(ctx, tx, contactusSpace.Record); err != nil {
		return response, err
	}

	happeningDto.UserIDs = params.Space.Data.UserIDs
	happeningDto.Status = "active"
	if happeningDto.Type == dbo4calendarium.HappeningTypeSingle {
		// TODO:add comment why we doing it only for single happening
		for _, slot := range happeningDto.Slots {
			if slot.Start.Date != "" {
				_ = happeningDto.AddDate(slot.Start.Date)
			}
			// TODO(help-wanted): populate dates between start & end dates
		}
	}

	contactsBySpaceID := make(map[string][]dal4contactus2.ContactEntry)

	//for participantID := range happeningDto.Participants {
	//	participantKey := dbmodels.SpaceItemID(participantID)
	//	spaceID := participantKey.Space()
	//	if spaceID == params.Space.ContactID {
	//		contactBrief := contactusSpace.Data.Contacts[participantKey.ItemID()]
	//		if contactBrief == nil {
	//			spaceContacts := contactsBySpaceID[teamID]
	//			if spaceContacts == nil {
	//				spaceContacts = make([]dal4contactus.DebtusSpaceContactEntry, 0, 1)
	//			}
	//			contactsBySpaceID[teamID] = append(teamContacts, dal4contactus.NewContactEntry(teamID, participantKey.ItemID()))
	//		} else {
	//			happeningDto.AddContact(teamID, participantKey.ItemID(), contactBrief)
	//		}
	//	} else {
	//		return response, errors.New("not implemented yet: adding participants from other teams at happening creation")
	//	}
	//}

	if len(contactsBySpaceID) > 0 {
		contactRecords := make([]dal.Record, 0)
		for _, teamContacts := range contactsBySpaceID {
			for _, contact := range teamContacts {
				contactRecords = append(contactRecords, contact.Record)
			}
		}
		if err = tx.GetMulti(ctx, contactRecords); err != nil {
			return response, err
		}
		//for teamID, teamContacts := range contactsBySpaceID {
		//	for _, contact := range teamContacts {
		//		happeningDto.AddContact(teamID, contact.ContactID, &contact.Data.ContactBrief)
		//	}
		//}
	}

	var happeningID string
	var happeningKey *dal.Key
	if happeningID, happeningKey, err = dal4spaceus2.GenerateNewSpaceModuleItemKey(
		ctx, tx, params.Space.ID, moduleID, happeningsCollection, 5, 10); err != nil {
		return response, err
	}
	response.ID = happeningID
	record := dal.NewRecordWithData(happeningKey, happeningDto)

	_ = dbo4linkage.UpdateRelatedIDs(params.Space.ID, &happeningDto.WithRelated, &happeningDto.WithRelatedIDs)

	if err = happeningDto.Validate(); err != nil {
		return response, fmt.Errorf("happening record is not valid for insertion: %w", err)
	}
	//panic("spaceDates: " + strings.Join(happeningDto.TeamDates, ","))
	if err = tx.Insert(ctx, record); err != nil {
		return response, fmt.Errorf("failed to insert new happening record: %w", err)
	}
	if happeningDto.Type == dbo4calendarium.HappeningTypeRecurring {
		params.SpaceModuleEntry.Record.MarkAsChanged()
		if params.SpaceModuleEntry.Data.RecurringHappenings == nil {
			params.SpaceModuleEntry.Data.RecurringHappenings = make(map[string]*dbo4calendarium.CalendarHappeningBrief)
		}
		params.SpaceModuleEntry.Data.RecurringHappenings[happeningID] = &dbo4calendarium.CalendarHappeningBrief{
			HappeningBrief: happeningDto.HappeningBrief,
			WithRelated:    happeningDto.WithRelated,
		}
		params.SpaceModuleUpdates = append(params.SpaceUpdates,
			update.ByFieldName("recurringHappenings."+happeningID, &happeningDto.HappeningBrief))
	}
	return
}
