package facade4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"testing"
)

// TODO: re-enable
func TestStartRetrospective(t *testing.T) {
	t.Skip("TODO: re-enable")

	//newRetrospectiveRef = func(teamRef *firestore.DocumentRef, id string) *firestore.DocumentRef {
	//	return &firestore.DocumentRef{InviteID: id, Path: "api4meetingus"}
	//}

	userContext := facade.NewUserContext("user1")
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
	const space1 = "space1"
	var validDurations = RetroDurations{Feedback: 2, Review: 5}

	var newRequest = func(id, spaceID string, durations RetroDurations) StartRetrospectiveRequest {
		return StartRetrospectiveRequest{
			RetroRequest: RetroRequest{
				MeetingID: id,
				SpaceRequest: dto4spaceus.SpaceRequest{
					SpaceID: spaceID,
				},
			},
			DurationsInMinutes: durations,
		}
	}

	t.Run("should_fail", func(t *testing.T) {
		t.Run("space_not_found", func(t *testing.T) {
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
			test(t, newRequest(UpcomingRetrospectiveID, space1, validDurations), expects{isNew: true})
		})
		t.Run("existing_record", func(t *testing.T) {
			test(t, newRequest(validRetroID, space1, validDurations), expects{isNew: false})
		})
	})
}
