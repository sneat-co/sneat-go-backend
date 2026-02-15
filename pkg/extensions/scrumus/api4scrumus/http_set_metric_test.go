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
	"github.com/stretchr/testify/assert"
)

func TestHttpPostSetMetric(t *testing.T) {
	setupMockVerify(t)

	oldSetMetric := setMetric
	t.Cleanup(func() { setMetric = oldSetMetric })

	t.Run("success", func(t *testing.T) {
		setMetric = func(ctx facade.ContextWithUser, request facade4scrumus.SetMetricRequest) (*facade4scrumus.SetMetricRequest, error) {
			val := 10
			return &facade4scrumus.SetMetricRequest{
				Request: facade4meetingus.Request{
					SpaceRequest: dto4spaceus.SpaceRequest{
						SpaceID: coretypes.SpaceID("s1"),
					},
					MeetingID: time.Now().Format("2006-01-02"),
				},
				Metric: "m1",
				MetricValue: dbo4scrumus.MetricValue{
					Int: &val,
				},
			}, nil
		}

		val := 10
		reqBody, _ := json.Marshal(facade4scrumus.SetMetricRequest{
			Request: facade4meetingus.Request{
				SpaceRequest: dto4spaceus.SpaceRequest{
					SpaceID: coretypes.SpaceID("s1"),
				},
				MeetingID: time.Now().Format("2006-01-02"),
			},
			Metric: "m1",
			MetricValue: dbo4scrumus.MetricValue{
				Int: &val,
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/set_metric", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		httpPostSetMetric(w, req)

		if w.Code == http.StatusInternalServerError {
			t.Errorf("expected success, got 500: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("bad_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v0/scrum/set_metric", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		httpPostSetMetric(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
