package api4scrumus

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestHttpDeleteTask(t *testing.T) {
	setupMockVerify(t)

	oldDeleteTask := deleteTask
	t.Cleanup(func() { deleteTask = oldDeleteTask })

	t.Run("success", func(t *testing.T) {
		deleteTask = func(ctx facade.ContextWithUser, request facade4scrumus.DeleteTaskRequest) (err error) {
			return nil
		}

		today := time.Now().Format("2006-01-02")
		req := httptest.NewRequest(http.MethodDelete, "/v0/scrum/delete_task?space=s1&date="+today+"&id=t1&type=todo&members=c1", nil)
		w := httptest.NewRecorder()

		httpDeleteTask(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("validation_error", func(t *testing.T) {
		// Missing space
		req := httptest.NewRequest(http.MethodDelete, "/v0/scrum/delete_task?date=2023-01-01&id=t1&type=todo", nil)
		w := httptest.NewRecorder()

		httpDeleteTask(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
