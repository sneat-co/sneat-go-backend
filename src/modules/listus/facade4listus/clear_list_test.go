package facade4listus

import (
	"context"
	"testing"
)

func TestClearList(t *testing.T) {
	// ClearList is currently empty, so we just test that it can be called.
	ClearList(context.Background(), nil, "test-list")
}
