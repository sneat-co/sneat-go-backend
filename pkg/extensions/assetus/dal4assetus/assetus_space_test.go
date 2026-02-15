package dal4assetus

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/const4assetus"
)

func TestAssetusRootKey(t *testing.T) {
	if AssetusRootKey == nil {
		t.Fatal("AssetusRootKey is nil")
	}
	if AssetusRootKey.Collection() != dbo4spaceus.SpaceModulesCollection {
		t.Errorf("expected collection %s, got %s", dbo4spaceus.SpaceModulesCollection, AssetusRootKey.Collection())
	}
	if AssetusRootKey.ID != const4assetus.ExtensionID {
		t.Errorf("expected ID %s, got %v", const4assetus.ExtensionID, AssetusRootKey.ID)
	}
}
