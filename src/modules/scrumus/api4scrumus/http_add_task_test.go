package api4scrumus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/facade4scrumus"
)

func TestAddTask(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Unexpected panic: %v", r)
		}
	}()
	// Call it with nils just to get coverage
	_, _ = addTask(nil, facade4scrumus.AddTaskRequest{})
}
