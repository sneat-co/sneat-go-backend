package dal4calendarium

import (
	"context"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/mocks/mock_dal"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"go.uber.org/mock/gomock"
)

func TestGetCalendariumSpace(t *testing.T) {
	ctx := context.Background()
	spaceID := coretypes.SpaceID("testspace")

	ctrl := gomock.NewController(t)
	tx := mock_dal.NewMockReadwriteTransaction(ctrl)

	var gotRecord dal.Record
	tx.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, record dal.Record) error {
		gotRecord = record
		return nil
	})

	entry, err := GetCalendariumSpace(ctx, tx, spaceID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if entry.Record == nil {
		t.Error("expected record to be set")
	}
	if gotRecord != entry.Record {
		t.Error("expected tx.Get to be called with entry.Record")
	}
}
