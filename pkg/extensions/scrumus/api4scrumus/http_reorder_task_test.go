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
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestHttpPostReorderTask(t *testing.T) {
	setupMockVerify(t)

	oldReorderTask := reorderTask
	t.Cleanup(func() { reorderTask = oldReorderTask })

	t.Run("success", func(t *testing.T) {
		reorderTask = func(ctx facade.ContextWithUser, request facade4scrumus.ReorderTaskRequest) (err error) {
			return nil
		}

		reqBody, _ := json.Marshal(facade4scrumus.ReorderTaskRequest{
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
			Len:   3,
			From:  0,
			To:    1,
			After: "t0",
		})
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/reorder_task", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		httpPostReorderTask(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/reorder_task", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		httpPostReorderTask(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
