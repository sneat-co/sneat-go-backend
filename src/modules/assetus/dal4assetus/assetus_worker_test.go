package dal4assetus

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
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
	t.Skip("skipping integration test for RunReadonlyAssetusSpaceWorker")
}

func TestRunAssetusSpaceWorkerTx(t *testing.T) {
	t.Skip("skipping integration test for RunAssetusSpaceWorkerTx")
}
