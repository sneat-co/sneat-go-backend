package analytics

import (
	//"github.com/strongo/measurement-protocol"
	"context"
	"testing"
)

func TestSendSingleMessage(t *testing.T) {
	var ctx context.Context
	if err := SendSingleMessage(ctx, nil); err == nil {
		t.Error("Expected to get error on nil context")
	}
}
