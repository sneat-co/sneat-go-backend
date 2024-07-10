package facade4scrumus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"testing"
	"time"
)

func TestDeleteTask(t *testing.T) {
	t.Skip("TODO: re-enable")
	//var db dal.DB
	//testdb.NewMockDB(t, db, testdb.WithProfile1())

	facade.GetDatabase = func(ctx context.Context) dal.DB {
		return nil //db
	}

	userContext := facade.NewUser("user1")

	ctx := context.Background()

	t.Run("empty_request", func(t *testing.T) {
		if err := DeleteTask(ctx, userContext, DeleteTaskRequest{}); err == nil {
			t.Fatal("Should fail on empty request")
		}
	})

	t.Run("valid_request", func(t *testing.T) {
		now := time.Now()
		request := DeleteTaskRequest{
			Request: facade4meetingus.Request{
				SpaceRequest: dto4teamus.SpaceRequest{
					SpaceID: "space1",
				},
				MeetingID: now.Format("2006-01-02"),
			},
			ContactID: "m1",
			Type:      "done",
			Task:      "d1",
		}

		t.Run("no_tasks", func(t *testing.T) {
			if err := DeleteTask(ctx, userContext, request); err != nil {
				t.Error(err)
			}
		})

		t.Run("existing_task", func(t *testing.T) {
			if err := DeleteTask(ctx, userContext, request); err != nil {
				t.Error(err)
			}
		})
	})

}
