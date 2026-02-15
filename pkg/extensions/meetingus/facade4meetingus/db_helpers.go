package facade4meetingus

import (
	"context"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	dal4contactus2 "github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// WorkerParams parameters
type WorkerParams struct {
	*dal4contactus2.ContactusSpaceWorkerParams
	Meeting WorkerMeeting
}

type workerItem struct {
	Key    *dal.Key
	Record dal.Record
}

func (v workerItem) GetID() string {
	return v.Key.ID.(string)
}

// WorkerSpaceDbo a worker for a space
type WorkerSpaceDbo struct {
	workerItem
}

// Data returns *dbo4spaceus.SpaceDbo
func (v WorkerSpaceDbo) Data() *dbo4spaceus.SpaceDbo {
	return v.Record.Data().(*dbo4spaceus.SpaceDbo)
}

// WorkerMeeting a worker for a meeting
type WorkerMeeting struct {
	workerItem
}

// Data returns *dbo4meetingus.Meeting
func (v WorkerMeeting) Data() *dbo4meetingus.Meeting {
	return v.Record.Data().(dbo4meetingus.MeetingInstance).BaseMeeting()
}

// RecordFactory a factory to create an api4meetingus record
type RecordFactory interface {
	// Collection name of collection
	Collection() string

	// NewRecordData creates an instance of api4meetingus record
	NewRecordData() dbo4meetingus.MeetingInstance
}

// Worker is a api4meetingus worker
type Worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params WorkerParams) (err error)

// RunMeetingWorker runs api4meetingus worker
func RunMeetingWorker(ctx context.Context, userCtx facade.UserContext, request Request, recordFactory RecordFactory, worker Worker) error {
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		params, err := GetMeetingAndSpace(ctx, tx, userCtx, request.SpaceID, request.MeetingID, recordFactory)
		if err != nil {
			return fmt.Errorf("failed to get api4meetingus & team records: %w", err)
		}
		return worker(ctx, tx, params)
	})
}

// GetMeetingAndSpace retrieve api4meetingus and team records
var GetMeetingAndSpace = func(ctx context.Context, tx dal.ReadwriteTransaction, userCtx facade.UserContext, spaceID coretypes.SpaceID, meetingID string, recordFactory RecordFactory) (params WorkerParams, err error) {
	params.ContactusSpaceWorkerParams = dal4contactus2.NewContactusSpaceWorkerParams(userCtx, spaceID)
	// Create team parameter
	// Create api4meetingus parameter
	meetingKey := dal.NewKeyWithParentAndID(params.Space.Key, recordFactory.Collection(), meetingID)
	params.Meeting = WorkerMeeting{
		workerItem: workerItem{
			Key:    params.Space.Key,
			Record: dal.NewRecordWithData(meetingKey, recordFactory.NewRecordData()),
		},
	}
	records := []dal.Record{
		params.Meeting.Record,
		params.Space.Record,
	}
	if err = tx.GetMulti(ctx, records); err != nil {
		return
	}

	if !params.Space.Record.Exists() {
		return params, fmt.Errorf("unknown team ContactID: %s", params.Space.Key.ID)
	}

	userBelongsToSpace := false
	teamData := params.Space.Data
	if len(teamData.UserIDs) == 0 {
		err = validation.NewErrBadRequestFieldValue("UserIDs",
			fmt.Sprintf("space record have no references to any user, key: %v; data: %+v",
				params.Space.Key.String(),
				teamData,
			))
		return
	}
	userID := userCtx.GetUserID()
	for _, v := range teamData.UserIDs {
		if v == userID {
			userBelongsToSpace = true
			break
		}
	}
	if !userBelongsToSpace {
		err = validation.NewErrBadRequestFieldValue("UserIDs", fmt.Sprintf("User does not belong to team, uid=%v, team.UserIDs: %+v", userID, teamData.UserIDs))
		return
	}

	if !params.Meeting.Record.Exists() {
		meeting := params.Meeting.Data()
		team := params.Space.Data
		contactusSpace := dal4contactus2.NewContactusSpaceEntry(params.Space.ID)
		if err := tx.Get(ctx, contactusSpace.Record); err != nil {
			return params, fmt.Errorf("failed to get contactus team record: %w", err)
		}
		meeting.UserIDs = team.UserIDs
		for contactID, teamMember := range contactusSpace.Data.Contacts {
			if teamMember.IsSpaceMember() {
				meeting.AddContact(spaceID, contactID, &dbo4meetingus.MeetingMemberBrief{ContactBrief: *teamMember})
			}
		}
	}
	return
}
