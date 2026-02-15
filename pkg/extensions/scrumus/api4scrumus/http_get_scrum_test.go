package api4scrumus

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestHttpGetScrum(t *testing.T) {
	setupMockVerify(t)

	oldGetScrum := getScrum
	t.Cleanup(func() { getScrum = oldGetScrum })

	t.Run("success", func(t *testing.T) {
		getScrum = func(ctx facade.ContextWithUser, request facade.IDRequest) (dbo4scrumus.Scrum, error) {
			return dbo4scrumus.Scrum{}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v0/scrum", nil)
		req.Header.Set("id", "scrum1")
		w := httptest.NewRecorder()

		httpGetScrum(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
