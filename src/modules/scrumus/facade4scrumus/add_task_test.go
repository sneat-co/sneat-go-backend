package facade4scrumus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-core/facade"
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {

	t.Skip("TODO: re-enable")
	//var db dal.DB
	//testdb.NewMockDB(t, db, testdb.WithProfile1())
	userContext := facade.NewUserContext("user1")

	t.Run("empty request", func(t *testing.T) {
		if _, err := AddTask(context.Background(), userContext, AddTaskRequest{}); err == nil {
			t.Fatal("should fail on empty request")
		}
	})

	t.Run("valid_requests", func(t *testing.T) {
		now := time.Now()
		request := AddTaskRequest{
			TaskRequest: TaskRequest{
				Request: facade4meetingus.Request{
					SpaceRequest: dto4spaceus.SpaceRequest{
						SpaceID: "space1",
					},
					MeetingID: now.Format("2006-01-02"),
				},
				ContactID: "m1",
				Type:      "done",
				Task:      "done1",
			},
			Title: "Test task",
		}

		t.Run("create_new_scrum", func(t *testing.T) {
			if _, err := AddTask(context.Background(), userContext, request); err != nil {
				t.Fatalf("should not fail on valid request, got: %v", err)
			}
		})

		//t.Run("update_existing_scrum", func(t *testing.T) {
		//	if _, err := AddTask(context.Background(), userContext, request); err != nil {
		//		t.Fatalf("should not fail on valid request, got: %v", err)
		//	}
		//})
	})
}
