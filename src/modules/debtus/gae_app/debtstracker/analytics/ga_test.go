package analytics

import (
	//"github.com/strongo/measurement-protocol"
	"context"
	"testing"
)

func TestSendSingleMessage(t *testing.T) {
	var c context.Context
	if err := SendSingleMessage(c, nil); err == nil {
		t.Error("Expected to get error on nil context")
	}
}
