package facade4retrospectus

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
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
					SpaceRequest: dto4spaceus.SpaceRequest{
						SpaceID: "space1",
					},
				},
				Item: "good1",
				From: dbo4retrospectus.TreePosition{Parent: "goods", Index: 0},
				To:   dbo4retrospectus.TreePosition{Parent: "goods", Index: 2},
			}

			ctx := context.Background()

			if err := MoveRetroItem(ctx, facade.NewUserContext("user1"), request); err == nil {
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
				SpaceRequest: dto4spaceus.SpaceRequest{
					SpaceID: "space1",
				},
			},
		}

		t.Run("moving within same parent", func(t *testing.T) {
			t.Skip("TODO: re-enable")
			request.Item = "g1"
			request.From = dbo4retrospectus.TreePosition{Parent: "goods", Index: 0}
			request.To = dbo4retrospectus.TreePosition{Parent: "goods", Index: 1}

			if err := MoveRetroItem(context.Background(), facade.NewUserContext("user1"), request); err != nil {
				t.Fatalf("failed to move retro item: %v", err)
			}
		})
	})
}
