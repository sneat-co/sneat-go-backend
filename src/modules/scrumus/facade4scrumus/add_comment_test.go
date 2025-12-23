package facade4scrumus

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func TestAddCommentRequest_Validate(t *testing.T) {
	request := AddCommentRequest{}
	err := request.Validate()
	if err == nil {
		t.Fatal("expected to get error on empty request")
	}
	request = AddCommentRequest{
		TaskRequest: TaskRequest{
			ContactID: "m1",
			Type:      "done",
			Task:      "task1",
			Request: facade4meetingus.Request{
				SpaceRequest: dto4spaceus.SpaceRequest{
					SpaceID: "space1",
				},
				MeetingID: "2020-12-13",
			},
		},
		Message: "message 1",
	}
	if err = request.Validate(); err != nil {
		t.Fatalf("unexpected error on valid request: %v", err)
	}
}

func TestAddComment(t *testing.T) {
	//userContext := facade4meetingus.NewUser("user1")

	t.Skip("TODO: re-enable")
	//var db dal.DB
	//testdb.NewMockDB(t, db, testdb.WithProfile1())

	facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
		return nil, nil //db
	}

	t.Run("add 1st comment", func(t *testing.T) {
		body := []byte(strings.ReplaceAll(strings.ReplaceAll(`{
	"spaceID":"space1",
	"meetingID":"2019-11-22",
	"type":"done",
	"task":"d1","
	members":"m1",
	"message":"msg1"
}`, "\n", ""), "\t", ""))
		var request AddCommentRequest
		if err := json.Unmarshal(body, &request); err != nil {
			t.Fatal(err)
		}

		//comment, err := AddComment(ctx, userContext, request)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//if comment.InviteID == "" {
		//	t.Error("InviteID is not set")
		//}
	})
}
