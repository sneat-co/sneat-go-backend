package dal4calendarium

import (
	"context"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

type mockTx struct {
	dal.ReadwriteTransaction
	gotRecord dal.Record
}

func (m *mockTx) Get(_ context.Context, record dal.Record) error {
	m.gotRecord = record
	return nil
}

func TestGetCalendariumSpace(t *testing.T) {
	ctx := context.Background()
	spaceID := coretypes.SpaceID("testspace")
	tx := &mockTx{}

	entry, err := GetCalendariumSpace(ctx, tx, spaceID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if entry.Record == nil {
		t.Error("expected record to be set")
	}
	if tx.gotRecord != entry.Record {
		t.Error("expected tx.Get to be called with entry.Record")
	}
}
