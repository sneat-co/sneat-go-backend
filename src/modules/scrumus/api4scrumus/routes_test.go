package api4scrumus

import (
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
