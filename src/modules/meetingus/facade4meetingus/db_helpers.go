package facade4meetingus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/models4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// WorkerParams parameters
type WorkerParams struct {
	*dal4contactus.ContactusTeamWorkerParams
	Meeting WorkerMeeting
}

type workerItem struct {
	Key    *dal.Key
	Record dal.Record
}

func (v workerItem) GetID() string {
	return v.Key.ID.(string)
}

// WorkerTeam DTO
type WorkerTeam struct {
	workerItem
}

// Data returns *models4teamus.TeamDbo
func (v WorkerTeam) Data() *models4teamus.TeamDbo {
	return v.Record.Data().(*models4teamus.TeamDbo)
}

// WorkerMeeting a worker for a meeting
type WorkerMeeting struct {
	workerItem
}

// Data returns *models4meetingus.Meeting
func (v WorkerMeeting) Data() *models4meetingus.Meeting {
	return v.Record.Data().(models4meetingus.MeetingInstance).BaseMeeting()
}

// RecordFactory a factory to create an api4meetingus record
type RecordFactory interface {
	// Collection name of collection
	Collection() string

	// NewRecordData creates an instance of api4meetingus record
	NewRecordData() models4meetingus.MeetingInstance
}

// Worker is a api4meetingus worker
type Worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params WorkerParams) (err error)

// RunMeetingWorker runs api4meetingus worker
func RunMeetingWorker(ctx context.Context, userID string, request Request, recordFactory RecordFactory, worker Worker) error {
	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		params, err := GetMeetingAndTeam(ctx, tx, userID, request.TeamID, request.MeetingID, recordFactory)
		if err != nil {
			return fmt.Errorf("failed to get api4meetingus & team records: %w", err)
		}
		return worker(ctx, tx, params)
	})
}

// GetMeetingAndTeam retrieve api4meetingus and team records
var GetMeetingAndTeam = func(ctx context.Context, tx dal.ReadwriteTransaction, uid, teamID, meetingID string, recordFactory RecordFactory) (params WorkerParams, err error) {
	params.ContactusTeamWorkerParams = dal4contactus.NewContactusTeamWorkerParams(uid, teamID)
	// Create team parameter
	// Create api4meetingus parameter
	meetingKey := dal.NewKeyWithParentAndID(params.Team.Key, recordFactory.Collection(), meetingID)
	params.Meeting = WorkerMeeting{
		workerItem: workerItem{
			Key:    params.Team.Key,
			Record: dal.NewRecordWithData(meetingKey, recordFactory.NewRecordData()),
		},
	}
	records := []dal.Record{
		params.Meeting.Record,
		params.Team.Record,
	}
	if err = tx.GetMulti(ctx, records); err != nil {
		return
	}

	if !params.Team.Record.Exists() {
		return params, fmt.Errorf("unknown team ContactID: %v", params.Team.Key.ID)
	}

	userBelongsToTeam := false
	teamData := params.Team.Data
	if len(teamData.UserIDs) == 0 {
		err = validation.NewErrBadRequestFieldValue("UserIDs",
			fmt.Sprintf("team record have no references to any user, key: %v; data: %+v",
				params.Team.Key.String(),
				teamData,
			))
		return
	}
	for _, v := range teamData.UserIDs {
		if v == uid {
			userBelongsToTeam = true
			break
		}
	}
	if !userBelongsToTeam {
		err = validation.NewErrBadRequestFieldValue("UserIDs", fmt.Sprintf("User does not belong to team, uid=%v, team.UserIDs: %+v", uid, teamData.UserIDs))
		return
	}

	if !params.Meeting.Record.Exists() {
		meeting := params.Meeting.Data()
		team := params.Team.Data
		contactusTeam := dal4contactus.NewContactusTeamModuleEntry(params.Team.ID)
		if err := tx.Get(ctx, contactusTeam.Record); err != nil {
			return params, fmt.Errorf("failed to get contactus team record: %w", err)
		}
		meeting.UserIDs = team.UserIDs
		for contactID, teamMember := range contactusTeam.Data.Contacts {
			if teamMember.IsTeamMember() {
				meeting.AddContact(teamID, contactID, &models4meetingus.MeetingMemberBrief{ContactBrief: *teamMember})
			}
		}
	}
	return
}
