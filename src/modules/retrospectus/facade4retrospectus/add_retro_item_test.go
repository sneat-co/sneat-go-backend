package facade4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"testing"
)

func TestAddRetroItem(t *testing.T) {

	//var db dal.DB
	//testdb.NewMockDB(t, db, testdb.WithProfile1())
	t.Skip("TODO: re-enable")

	userContext := facade.NewUser("user1")
	t.Run("should_succeed", func(t *testing.T) {
		t.Run("upcoming_retrospective", func(t *testing.T) {
			newSpaceKey = func(id string) *dal.Key {
				return dal.NewKeyWithID(dal4teamus.SpacesCollection, id)
			}

			request := AddRetroItemRequest{
				RetroItemRequest: RetroItemRequest{
					Request: facade4meetingus.Request{
						SpaceRequest: dto4teamus.SpaceRequest{
							SpaceID: "space1",
						},
						MeetingID: UpcomingRetrospectiveID,
					},
					Type: "good",
				},
				Title: "Good # 1",
			}

			ctx := context.Background()

			_, _ = AddRetroItem(ctx, userContext, request)
			//if _, _ = AddRetroItem(ctx, userContext, request); false {
			// TODO: t.Fatalf("failed to add retro item: %v", err)
			//}
		})
	})
}
