package api4retrospectus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/facade4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/apicore/httpmock"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddRetroItem(t *testing.T) {
	called := 0
	addRetroItem = func(_ context.Context, _ facade.User, request facade4retrospectus.AddRetroItemRequest) (response facade4retrospectus.AddRetroItemResponse, err error) {
		called++
		return response, err
	}
	req := httpmock.NewPostJSONRequest(http.MethodPost, "/v0/retrospective/add_retro_item", &facade4retrospectus.AddRetroItemRequest{
		RetroItemRequest: facade4retrospectus.RetroItemRequest{
			Request: facade4meetingus.Request{
				SpaceRequest: dto4spaceus.SpaceRequest{
					SpaceID: "space1",
				},
				MeetingID: "retro1",
			},
			Type: "good",
		},
		Title: "New item #1",
	})
	w := httptest.NewRecorder()

	verifyRequest = func(w http.ResponseWriter, r *http.Request, options verify.RequestOptions) (ctx context.Context, userContext facade.User, err error) {
		return r.Context(), nil, nil
	}

	handler := http.HandlerFunc(httpPostAddRetroItem)
	handler.ServeHTTP(w, req)
	responseBody := w.Body.String()

	if w.Code != http.StatusCreated {
		t.Fatalf("expected to get status code %v, got %v; response body: %v",
			http.StatusCreated, w.Code, responseBody)
	}
	switch called {
	case 0:
		t.Errorf("addRetroItem have not been called: %v", responseBody)
	case 1:
		break
	default:
		t.Errorf("addRetroItem expetect to be called just once, was called %v times", called)
	}
}
