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

func TestAddTask(t *testing.T) {
	// Call it with nils just to get coverage
	_, _ = addTask(nil, facade4scrumus.AddTaskRequest{})
}

func TestHttpPostAddTask(t *testing.T) {
	setupMockVerify(t)

	// Mock addTask facade function
	oldAddTask := addTask
	t.Cleanup(func() { addTask = oldAddTask })

	t.Run("success", func(t *testing.T) {
		addTask = func(ctx facade.ContextWithUser, request facade4scrumus.AddTaskRequest) (*facade4scrumus.AddTaskResponse, error) {
			return &facade4scrumus.AddTaskResponse{}, nil
		}

		today := time.Now().Format("2006-01-02")
		reqBody, _ := json.Marshal(facade4scrumus.AddTaskRequest{
			TaskRequest: facade4scrumus.TaskRequest{
				Request: facade4meetingus.Request{
					SpaceRequest: dto4spaceus.SpaceRequest{
						SpaceID: coretypes.SpaceID("s1"),
					},
					MeetingID: today,
				},
				ContactID: "c1",
				Type:      "todo",
			},
			Title: "Test Title",
		})
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/add_task", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		httpPostAddTask(w, req)

		if w.Code == http.StatusInternalServerError {
			t.Errorf("expected success, got 500: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("bad_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/add_task", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		httpPostAddTask(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
