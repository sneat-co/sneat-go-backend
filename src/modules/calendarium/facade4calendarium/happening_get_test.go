package facade4calendarium

import (
	"context"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

type mockTx struct {
	dal.ReadwriteTransaction
	getFunc func(ctx context.Context, record dal.Record) error
}

func (m *mockTx) Get(ctx context.Context, record dal.Record) error {
	return m.getFunc(ctx, record)
}

func TestGetForUpdate(t *testing.T) {
	ctx := context.Background()
	spaceID := coretypes.SpaceID("test_space")
	happeningID := "test_happening"
	dto := dbo4calendarium.HappeningDbo{}

	called := false
	mock := &mockTx{
		getFunc: func(ctx context.Context, record dal.Record) error {
			called = true
			if record.Key().ID != happeningID {
				t.Errorf("expected happeningID %s, got %v", happeningID, record.Key().ID)
			}
			return nil
		},
	}

	record, err := GetForUpdate(ctx, mock, spaceID, happeningID, dto)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected Get to be called")
	}
	if record == nil {
		t.Fatal("expected record, got nil")
	}
}
