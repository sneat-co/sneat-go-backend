package reminders

import (
	"net/http/httptest"
	"testing"
)

func TestAllowOrigin(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	allowOrigin(responseRecorder)
	header := responseRecorder.Header()
	if v := header.Get("Access-Control-Allow-Origin"); v != "*" {
		t.Errorf("Expected to get '*', got: %v", v)
	}
}
