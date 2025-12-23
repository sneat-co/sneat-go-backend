package dal4assetus

import (
	"context"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

func TestNewAssetusSpaceWorkerParams(t *testing.T) {
	spaceID := coretypes.SpaceID("testspace")
	params := NewAssetusSpaceWorkerParams(nil, spaceID)
	if params == nil {
		t.Fatal("expected params, got nil")
	}
	if params.Space.ID != spaceID {
		t.Errorf("expected spaceID %s, got %s", spaceID, params.Space.ID)
	}
}

func TestRunReadonlyAssetusSpaceWorker(t *testing.T) {
	orig := runReadonlyModuleSpaceWorker
	defer func() { runReadonlyModuleSpaceWorker = orig }()

	called := false
	runReadonlyModuleSpaceWorker = func(
		ctx context.Context,
		userCtx facade.UserContext,
		request dto4spaceus.SpaceRequest,
		moduleID coretypes.ExtID,
		data *dbo4assetus.AssetusSpaceDbo,
		worker func(ctx context.Context, tx dal.ReadTransaction, spaceWorkerParams *AssetusSpaceWorkerParams) (err error),
	) error {
		called = true
		if moduleID != const4assetus.ExtensionID {
			t.Errorf("expected moduleID %s, got %s", const4assetus.ExtensionID, moduleID)
		}
		return nil
	}

	err := RunReadonlyAssetusSpaceWorker(context.Background(), nil, dto4spaceus.SpaceRequest{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("runReadonlyModuleSpaceWorker was not called")
	}
}

func TestRunAssetusSpaceWorkerTx(t *testing.T) {
	orig := runModuleSpaceWorkerTx
	defer func() { runModuleSpaceWorkerTx = orig }()

	called := false
	runModuleSpaceWorkerTx = func(
		ctx facade.ContextWithUser,
		tx dal.ReadwriteTransaction,
		spaceID coretypes.SpaceID,
		moduleID coretypes.ExtID,
		data *dbo4assetus.AssetusSpaceDbo,
		worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, spaceWorkerParams *AssetusSpaceWorkerParams) (err error),
	) error {
		called = true
		if moduleID != const4assetus.ExtensionID {
			t.Errorf("expected moduleID %s, got %s", const4assetus.ExtensionID, moduleID)
		}
		return nil
	}

	err := RunAssetusSpaceWorkerTx(nil, nil, "test-space", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("runModuleSpaceWorkerTx was not called")
	}
}
