package api4scrumus

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func TestHttpPostAddComment(t *testing.T) {
	setupMockVerify(t)

	oldAddComment := addComment
	t.Cleanup(func() { addComment = oldAddComment })

	t.Run("success", func(t *testing.T) {
		addComment = func(ctx facade.ContextWithUser, request facade4scrumus.AddCommentRequest) (*dbo4scrumus.Comment, error) {
			return &dbo4scrumus.Comment{
				ID:      "c1",
				Message: "test",
				By: &dbmodels.ByUser{
					UID:   "u1",
					Title: "User 1",
				},
			}, nil
		}

		reqBody, _ := json.Marshal(facade4scrumus.AddCommentRequest{
			TaskRequest: facade4scrumus.TaskRequest{
				Request: facade4meetingus.Request{
					SpaceRequest: dto4spaceus.SpaceRequest{
						SpaceID: coretypes.SpaceID("s1"),
					},
					MeetingID: time.Now().Format("2006-01-02"),
				},
				ContactID: "c1",
				Type:      "todo",
				Task:      "t1",
			},
			Message: "Test comment",
		})
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/add_comment", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		httpPostAddComment(w, req)

		if w.Code == http.StatusInternalServerError {
			t.Errorf("expected success, got 500: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("bad_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/add_comment", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		httpPostAddComment(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
