package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"strings"
)

// CreateHappening creates a recurring happening
func CreateHappening(
	ctx context.Context, user facade.User, request dto4calendarium.CreateHappeningRequest,
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
				CreatedBy: user.GetID(),
			},
		},
		//WithTeamDates: dbmodels.WithTeamDates{
		//	WithTeamIDs: dbmodels.WithTeamIDs{
		//		TeamIDs: []string{request.TeamID},
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
	err = dal4teamus.CreateTeamItem(ctx, user, request.TeamRequest,
		const4calendarium.ModuleID,
		new(dbo4calendarium.CalendariumTeamDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo]) (err error) {
			response, err = createHappeningTx(ctx, tx, happeningDto, params)
			return
		},
	)
	response.Dto = *happeningDto
	return
}

func createHappeningTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	happeningDto *dbo4calendarium.HappeningDbo,
	params *dal4teamus.ModuleTeamWorkerParams[*dbo4calendarium.CalendariumTeamDbo],
) (
	response dto4calendarium.CreateHappeningResponse, err error,
) {
	happeningDto.CreatedAt = params.Started
	contactusTeam := dal4contactus.NewContactusTeamModuleEntry(params.Team.ID)
	if err = params.GetRecords(ctx, tx, contactusTeam.Record); err != nil {
		return response, err
	}

	happeningDto.UserIDs = params.Team.Data.UserIDs
	happeningDto.Status = "active"
	if happeningDto.Type == dbo4calendarium.HappeningTypeSingle {
		for _, slot := range happeningDto.Slots {
			if slot.Start.Date != "" {
				happeningDto.Dates = append(happeningDto.Dates, slot.Start.Date)
			}
			if happeningDto.DateMin == "" || happeningDto.DateMin > slot.Start.Date {
				happeningDto.DateMin = slot.Start.Date
			}
			endDate := slot.End.Date
			if endDate == "" {
				endDate = slot.Start.Date
			}
			if happeningDto.DateMax == "" && endDate != "" && happeningDto.DateMax < endDate {
				happeningDto.DateMax = endDate
			}
			// TODO(help-wanted): populate dates between start & end dates
		}
	}

	contactsByTeamID := make(map[string][]dal4contactus.ContactEntry)

	//for participantID := range happeningDto.Participants {
	//	participantKey := dbmodels.TeamItemID(participantID)
	//	teamID := participantKey.TeamID()
	//	if teamID == params.Team.ID {
	//		contactBrief := contactusTeam.Data.Contacts[participantKey.ItemID()]
	//		if contactBrief == nil {
	//			teamContacts := contactsByTeamID[teamID]
	//			if teamContacts == nil {
	//				teamContacts = make([]dal4contactus.ContactEntry, 0, 1)
	//			}
	//			contactsByTeamID[teamID] = append(teamContacts, dal4contactus.NewContactEntry(teamID, participantKey.ItemID()))
	//		} else {
	//			happeningDto.AddContact(teamID, participantKey.ItemID(), contactBrief)
	//		}
	//	} else {
	//		return response, errors.New("not implemented yet: adding participants from other teams at happening creation")
	//	}
	//}

	if len(contactsByTeamID) > 0 {
		contactRecords := make([]dal.Record, 0)
		for _, teamContacts := range contactsByTeamID {
			for _, contact := range teamContacts {
				contactRecords = append(contactRecords, contact.Record)
			}
		}
		if err = tx.GetMulti(ctx, contactRecords); err != nil {
			return response, err
		}
		//for teamID, teamContacts := range contactsByTeamID {
		//	for _, contact := range teamContacts {
		//		happeningDto.AddContact(teamID, contact.ID, &contact.Data.ContactBrief)
		//	}
		//}
	}

	var happeningID string
	var happeningKey *dal.Key
	if happeningID, happeningKey, err = dal4teamus.GenerateNewTeamModuleItemKey(
		ctx, tx, params.Team.ID, moduleID, happeningsCollection, 5, 10); err != nil {
		return response, err
	}
	response.ID = happeningID
	record := dal.NewRecordWithData(happeningKey, happeningDto)

	_ = dbo4linkage.UpdateRelatedIDs(&happeningDto.WithRelated, &happeningDto.WithRelatedIDs)

	if err = happeningDto.Validate(); err != nil {
		return response, fmt.Errorf("happening record is not valid for insertion: %w", err)
	}
	//panic("teamDates: " + strings.Join(happeningDto.TeamDates, ","))
	if err = tx.Insert(ctx, record); err != nil {
		return response, fmt.Errorf("failed to insert new happening record: %w", err)
	}
	if happeningDto.Type == dbo4calendarium.HappeningTypeRecurring {
		if params.TeamModuleEntry.Data.RecurringHappenings == nil {
			params.TeamModuleEntry.Data.RecurringHappenings = make(map[string]*dbo4calendarium.CalendarHappeningBrief)
		}
		params.TeamModuleEntry.Data.RecurringHappenings[happeningID] = &dbo4calendarium.CalendarHappeningBrief{
			HappeningBrief: happeningDto.HappeningBrief,
		}
		params.TeamModuleUpdates = append(params.TeamUpdates, dal.Update{
			Field: "recurringHappenings." + happeningID,
			Value: &happeningDto.HappeningBrief,
		})
	}
	return
}
