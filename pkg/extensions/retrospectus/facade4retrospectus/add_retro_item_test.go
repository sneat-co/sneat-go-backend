package facade4retrospectus

import (
	"context"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

func TestAddRetroItem(t *testing.T) {

	//var db dal.DB
	//testdb.NewMockDB(t, db, testdb.WithProfile1())
	t.Skip("TODO: re-enable")

	userContext := facade.NewUserContext("user1")
	t.Run("should_succeed", func(t *testing.T) {
		t.Run("upcoming_retrospective", func(t *testing.T) {
			newSpaceKey = func(id coretypes.SpaceID) *dal.Key {
				return dal.NewKeyWithID(dbo4spaceus.SpacesCollection, id)
			}

			request := AddRetroItemRequest{
				RetroItemRequest: RetroItemRequest{
					Request: facade4meetingus.Request{
						SpaceRequest: dto4spaceus.SpaceRequest{
							SpaceID: "space1",
						},
						MeetingID: UpcomingRetrospectiveID,
					},
					Type: "good",
				},
				Title: "Good # 1",
			}

			ctx := facade.NewContextWithUserID(context.Background(), userContext.GetUserID())

			_, _ = AddRetroItem(ctx, request)
			//if _, _ = AddRetroItem(ctx, userContext, request); false {
			// TODO: t.Fatalf("failed to add retro item: %v", err)
			//}
		})
	})
}
