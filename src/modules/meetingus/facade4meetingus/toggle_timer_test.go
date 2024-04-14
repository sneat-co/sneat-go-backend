package facade4meetingus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/models4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"testing"
	"time"
)

type recordFactory struct {
}

// Collection "api4meetingus"
func (recordFactory) Collection() string {
	return "api4meetingus"
}

// NewRecord creates new record
func (recordFactory) NewRecordData() models4meetingus.MeetingInstance {
	return &models4meetingus.Meeting{}
}

func TestToggleTimer(t *testing.T) { // TODO(help-wanted): add more test cases
	t.Skip("TODO: re-enable")
	//var db dal.DB
	//testdb.NewMockDB(t, db, testdb.WithProfile1())

	userContext := facade.NewUser("user1")

	const (
		team1ID = "team1"
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
		initMeeting func(meeting *models4meetingus.Meeting),
		assert func(response ToggleTimerResponse, meeting models4meetingus.Meeting, team models4teamus.TeamDbo),
	) {
		assertTimer := func(source string, timer *models4meetingus.Timer) {
			if timer == nil {
				t.Fatalf(source + ".Timer == nil")
			}
			if timer.Status != expected.status {
				t.Errorf(source+".Timer.Status != expected.status: `%v` != `%v`", timer.Status, expected.status)
			}
			if timer.At.IsZero() {
				t.Error(source + ".Timer.At is zero")
			}
		}

		ctx := context.Background()
		response, err := ToggleTimer(ctx, userContext, ToggleParams{Params: Params{recordFactory{}, nil}, Request: request})
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
				TeamRequest: dto4teamus.TeamRequest{
					TeamID: team1ID,
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
					func(meeting *models4meetingus.Meeting) {
						if meeting == nil {
							panic("required parameter 'api4meetingus *models4meetingus.MeetingID' is nil")
						}
						now := time.Now()
						meeting.Started = &now
						meeting.Version = 1
						meeting.Timer = &models4meetingus.Timer{
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
					func(meeting *models4meetingus.Meeting) {

					},
					func(response ToggleTimerResponse, meeting models4meetingus.Meeting, team models4teamus.TeamDbo) {
						if meeting.Timer.ActiveMemberID != request.Member {
							t.Errorf("api4meetingus.Timer.ActiveMemberID !== request.MemberDto: %v != %v", meeting.Timer.ActiveMemberID, request.Member)
						}
					})
			})
		})
	})
}
