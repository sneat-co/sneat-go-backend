package queues

import "testing"

func TestKnownQueues(t *testing.T) {
	if len(KnownQueues) == 0 {
		t.Errorf("KnownQueues is empty")
	}
}
