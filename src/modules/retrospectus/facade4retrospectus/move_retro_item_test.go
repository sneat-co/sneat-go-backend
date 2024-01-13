package facade4retrospectus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/models4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"testing"
)

func TestMoveRetroItem(t *testing.T) {
	t.Skip("TODO: re-enable")
	//var db dal.DB
	//testdb.NewMockDB(t, db, testdb.WithProfile1())
	//const uid = "123"

	t.Run("Should fail", func(t *testing.T) {
		t.Run("when retrospective not found by ContactID", func(t *testing.T) {
			request := MoveRetroItemRequest{
				Request: RetroRequest{
					MeetingID: "non_existing_retro",
					TeamRequest: dto4teamus.TeamRequest{
						TeamID: "team1",
					},
				},
				Item: "good1",
				From: models4retrospectus.TreePosition{Parent: "goods", Index: 0},
				To:   models4retrospectus.TreePosition{Parent: "goods", Index: 2},
			}

			ctx := context.Background()

			if err := MoveRetroItem(ctx, facade.NewUser("user1"), request); err == nil {
				t.Fatal("Should fail")
			} else {
				t.Logf("Failed as expects: %v", err)
			}
		})
	})

	t.Run("Should succeed", func(t *testing.T) {
		request := MoveRetroItemRequest{
			Request: RetroRequest{
				MeetingID: "retro1",
				TeamRequest: dto4teamus.TeamRequest{
					TeamID: "team1",
				},
			},
		}

		t.Run("moving within same parent", func(t *testing.T) {
			t.Skip("TODO: re-enable")
			request.Item = "g1"
			request.From = models4retrospectus.TreePosition{Parent: "goods", Index: 0}
			request.To = models4retrospectus.TreePosition{Parent: "goods", Index: 1}

			if err := MoveRetroItem(context.Background(), facade.NewUser("user1"), request); err != nil {
				t.Fatalf("failed to move retro item: %v", err)
			}
		})
	})
}
