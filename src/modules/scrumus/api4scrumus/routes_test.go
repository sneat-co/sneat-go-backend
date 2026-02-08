package api4scrumus

import (
	"net/http"
	"testing"
)

func TestScrumAPI(t *testing.T) {
	if getScrum == nil {
		t.Error("getScrum is nil")
	}
	if thumbUp == nil {
		t.Error("thumbUp is nil")
	}
}

func TestRegisterHttpRoutes(t *testing.T) {
	routes := make(map[string]bool)
	handle := func(method, path string, handler http.HandlerFunc) {
		routes[method+":"+path] = true
	}
	RegisterHttpRoutes(handle)

	expectedRoutes := []string{
		"GET:/v0/scrum",
		"POST:/v0/scrum/add_task",
		"POST:/v0/scrum/set_metric",
		"POST:/v0/scrum/reorder_task",
		"POST:/v0/scrum/add_comment",
		"DELETE:/v0/scrum/delete_task",
		"POST:/v0/scrum/thumb_up_task",
		"POST:/v0/scrum/toggle_meeting_timer",
		"POST:/v0/scrum/toggle_member_timer",
	}

	for _, route := range expectedRoutes {
		if !routes[route] {
			t.Errorf("Expected route %s not registered", route)
		}
	}
}
