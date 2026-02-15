package facade4calendarium

import (
	"context"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/mocks/mock_dal"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"go.uber.org/mock/gomock"
)

func TestGetForUpdate(t *testing.T) {
	ctx := context.Background()
	spaceID := coretypes.SpaceID("test_space")
	happeningID := "test_happening"
	dto := dbo4calendarium.HappeningDbo{}

	ctrl := gomock.NewController(t)
	tx := mock_dal.NewMockReadwriteTransaction(ctrl)

	called := false
	tx.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, record dal.Record) error {
		called = true
		if record.Key().ID != happeningID {
			t.Errorf("expected happeningID %s, got %v", happeningID, record.Key().ID)
		}
		return nil
	})

	record, err := GetForUpdate(ctx, tx, spaceID, happeningID, dto)
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
