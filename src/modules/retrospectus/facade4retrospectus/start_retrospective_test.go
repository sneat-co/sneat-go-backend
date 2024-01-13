package facade4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"testing"
)

// TODO: re-enable
func TestStartRetrospective(t *testing.T) {
	t.Skip("TODO: re-enable")

	//newRetrospectiveRef = func(teamRef *firestore.DocumentRef, id string) *firestore.DocumentRef {
	//	return &firestore.DocumentRef{InviteID: id, Path: "api4meetingus"}
	//}

	userContext := facade.NewUser("user1")
	ctx := context.Background()

	type expects struct {
		isNew bool
	}

	test := func(t *testing.T, request StartRetrospectiveRequest, expected expects) {
		// SUT call
		response, isNew, err := StartRetrospective(ctx, userContext, request)

		if err != nil {
			t.Fatalf("Unexptected error: %v", err)
		}
		if isNew != expected.isNew {
			if expected.isNew {
				t.Errorf("expected to be a new record, response: %+v", response)
			} else {
				t.Fatalf("expected to be an existing record")
			}
		}
		if response.ID == "" {
			t.Errorf("response.InviteID is empty string")
		}
		if response.Data == nil {
			t.Fatalf("response.MeetingRecord == nil")
		}
	}

	const validRetroID = "retro1"
	const team1 = "team1"
	var validDurations = RetroDurations{Feedback: 2, Review: 5}

	var newRequest = func(id, teamID string, durations RetroDurations) StartRetrospectiveRequest {
		return StartRetrospectiveRequest{
			RetroRequest: RetroRequest{
				MeetingID: id,
				TeamRequest: dto4teamus.TeamRequest{
					TeamID: teamID,
				},
			},
			DurationsInMinutes: durations,
		}
	}

	t.Run("should_fail", func(t *testing.T) {
		t.Run("team_not_found", func(t *testing.T) {
			if _, _, err := StartRetrospective(ctx, userContext, newRequest(validRetroID, "invalidteamid", validDurations)); err == nil {
				t.Fatal("expected to get error")
			}
		})
	})

	t.Run("should_succeed", func(t *testing.T) {
		t.Run("new_upcoming_record", func(t *testing.T) {
			txGetRetrospective = func(ctx context.Context, tx dal.ReadwriteTransaction, record dal.Record) (err error) {
				return nil
			}
			test(t, newRequest(UpcomingRetrospectiveID, team1, validDurations), expects{isNew: true})
		})
		t.Run("existing_record", func(t *testing.T) {
			test(t, newRequest(validRetroID, team1, validDurations), expects{isNew: false})
		})
	})
}
