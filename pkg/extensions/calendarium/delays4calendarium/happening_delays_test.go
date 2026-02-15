package delays4calendarium

import (
	"context"
	"testing"

	"github.com/strongo/delaying"
)

type mockDelayer struct {
	delaying.Delayer
	enqueued bool
}

func (m *mockDelayer) EnqueueWork(_ context.Context, params delaying.Params, args ...any) error {
	_, _ = params, args
	m.enqueued = true
	return nil
}

func TestDelayUpdateHappeningBrief(t *testing.T) {
	mock := &mockDelayer{}
	InitDelaying(func(key string, i any) delaying.Delayer {
		return mock
	})

	err := DelayUpdateHappeningBrief(context.Background(), "user1", "space1", "happening1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !mock.enqueued {
		t.Error("expected work to be enqueued")
	}
}
