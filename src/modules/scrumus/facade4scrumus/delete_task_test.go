package facade4scrumus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-core/facade"
	"testing"
	"time"
)

func TestDeleteTask(t *testing.T) {
	t.Skip("TODO: re-enable")
	//var db dal.DB
	//testdb.NewMockDB(t, db, testdb.WithProfile1())

	facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
		return nil, nil //db
	}

	ctx := facade.NewContextWithUser(context.Background(), "user1")

	t.Run("empty_request", func(t *testing.T) {
		if err := DeleteTask(ctx, DeleteTaskRequest{}); err == nil {
			t.Fatal("Should fail on empty request")
		}
	})

	t.Run("valid_request", func(t *testing.T) {
		now := time.Now()
		request := DeleteTaskRequest{
			Request: facade4meetingus.Request{
				SpaceRequest: dto4spaceus.SpaceRequest{
					SpaceID: "space1",
				},
				MeetingID: now.Format("2006-01-02"),
			},
			ContactID: "m1",
			Type:      "done",
			Task:      "d1",
		}

		t.Run("no_tasks", func(t *testing.T) {
			if err := DeleteTask(ctx, request); err != nil {
				t.Error(err)
			}
		})

		t.Run("existing_task", func(t *testing.T) {
			if err := DeleteTask(ctx, request); err != nil {
				t.Error(err)
			}
		})
	})

}
