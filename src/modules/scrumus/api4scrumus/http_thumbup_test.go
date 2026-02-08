package api4scrumus

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestHttpPostThumbUp(t *testing.T) {
	setupMockVerify(t)

	oldThumbUp := thumbUp
	t.Cleanup(func() { thumbUp = oldThumbUp })

	t.Run("success", func(t *testing.T) {
		thumbUp = func(ctx facade.ContextWithUser, request facade4scrumus.ThumbUpRequest) (err error) {
			return nil
		}

		reqBody, _ := json.Marshal(facade4scrumus.ThumbUpRequest{
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
			Value: true,
		})
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/thumb_up_task", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		httpPostThumbUp(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/thumb_up_task", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		httpPostThumbUp(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
