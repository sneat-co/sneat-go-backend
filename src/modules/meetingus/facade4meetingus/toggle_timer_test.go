package facade4meetingus

import (
	"context"
	"testing"
	"time"

	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

type recordFactory struct {
}

// Collection "api4meetingus"
func (recordFactory) Collection() string {
	return "api4meetingus"
}

// NewRecord creates new record
func (recordFactory) NewRecordData() dbo4meetingus.MeetingInstance {
	return &dbo4meetingus.Meeting{}
}

func TestToggleTimerRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     ToggleTimerRequest
		wantErr bool
	}{
		{"valid", ToggleTimerRequest{
			Operation: TimerOpStart,
			Request: Request{
				SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: "s1"},
				MeetingID:    "m1",
			},
		}, false},
		{"missing_op", ToggleTimerRequest{
			Request: Request{
				SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: "s1"},
				MeetingID:    "m1",
			},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ToggleTimerRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToggleTimer(t *testing.T) { // TODO(help-wanted): add more test cases
	t.Skip("TODO: re-enable")
	//var db dal.DB
	//testdb.NewMockDB(t, db, testdb.WithProfile1())

	const (
		space1ID = "space1"
	)

	type expecting struct {
		status string
	}

	testToggleTimer := func(
		t *testing.T,
		existingMeetingRecord bool,
		request ToggleTimerRequest,
		timestamps []dbmodels.Timestamp,
		expected expecting,
		initMeeting func(meeting *dbo4meetingus.Meeting),
		assert func(response ToggleTimerResponse, meeting dbo4meetingus.Meeting, team dbo4spaceus.SpaceDbo),
	) {
		assertTimer := func(source string, timer *dbo4meetingus.Timer) {
			if timer == nil {
				t.Fatal(source + ".Timer == nil")
			}
			if timer.Status != expected.status {
				t.Errorf(source+".Timer.Status != expected.status: `%s` != `%s`", timer.Status, expected.status)
			}
			if timer.At.IsZero() {
				t.Error(source + ".Timer.At is zero")
			}
		}

		ctx := facade.NewContextWithUserID(context.Background(), "user1")
		response, err := ToggleTimer(ctx, ToggleParams{Params: Params{recordFactory{}, nil}, Request: request})
		if err != nil {
			t.Fatal(err)
		}

		assertTimer("response", response.Timer)
		//assertTimer("api4meetingus", meeting.Timer)
		//if assert != nil {
		//	assert(response, *meeting, team)
		//}
	}

	newRequest := func(op string, member string) ToggleTimerRequest {
		return ToggleTimerRequest{
			Operation: op,
			Member:    member,
			Request: Request{
				SpaceRequest: dto4spaceus.SpaceRequest{
					SpaceID: space1ID,
				},
				MeetingID: "2010-11-22",
			},
		}
	}

	t.Run("toggle_meeting_timer", func(t *testing.T) {

		t.Run("existing_record", func(t *testing.T) {

			t.Run("first_start", func(t *testing.T) {
				request := newRequest(TimerOpStart, "")
				testToggleTimer(t, true, request, nil,
					expecting{status: TimerStatusActive},
					nil,
					nil,
				)
			})

			t.Run("pause", func(t *testing.T) {
				request := newRequest(TimerOpPause, "")
				testToggleTimer(t, true, request, nil,
					expecting{status: TimerStatusPaused},
					func(meeting *dbo4meetingus.Meeting) {
						if meeting == nil {
							panic("required parameter 'api4meetingus *dbo4meetingus.MeetingID' is nil")
						}
						now := time.Now()
						meeting.Started = &now
						meeting.Version = 1
						meeting.Timer = &dbo4meetingus.Timer{
							By: dbmodels.ByUser{
								UID: "u1",
							},
							Status: TimerStatusActive,
							At:     now,
						}
					},
					nil,
				)
			})

		})

	})

	t.Run("toggle_member_timer", func(t *testing.T) {
		t.Run("existing_record", func(t *testing.T) {
			t.Run("first_start", func(t *testing.T) {
				request := newRequest("start", "m1")
				testToggleTimer(t, true, request, nil,
					expecting{status: TimerStatusActive},
					func(meeting *dbo4meetingus.Meeting) {

					},
					func(response ToggleTimerResponse, meeting dbo4meetingus.Meeting, team dbo4spaceus.SpaceDbo) {
						if meeting.Timer.ActiveMemberID != request.Member {
							t.Errorf("api4meetingus.Timer.ActiveMemberID !== request.MemberDto: %v != %v", meeting.Timer.ActiveMemberID, request.Member)
						}
					})
			})
		})
	})
}
