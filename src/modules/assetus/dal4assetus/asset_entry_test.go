package dal4assetus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestNewAssetKeyWithoutID(t *testing.T) {
	spaceID := coretypes.SpaceID("testspace")
	key := NewAssetKeyWithoutID(spaceID)
	if key == nil {
		t.Fatal("expected key, got nil")
	}
	if key.ID != nil {
		t.Errorf("expected nil ID for incomplete key, got %v", key.ID)
	}
	if key.Collection() != dbo4assetus.SpaceAssetsCollection {
		t.Errorf("expected collection %s, got %s", dbo4assetus.SpaceAssetsCollection, key.Collection())
	}
	if key.Parent() == nil {
		t.Fatal("expected parent key, got nil")
	}
}

func TestNewAssetKey(t *testing.T) {
	spaceID := coretypes.SpaceID("testspace")
	assetID := "test_asset_123"
	key := NewAssetKey(spaceID, assetID)
	if key == nil {
		t.Fatal("expected key, got nil")
	}
	if key.ID != assetID {
		t.Errorf("expected ID %s, got %v", assetID, key.ID)
	}
	if key.Collection() != dbo4assetus.SpaceAssetsCollection {
		t.Errorf("expected collection %s, got %s", dbo4assetus.SpaceAssetsCollection, key.Collection())
	}
	if key.Parent() == nil {
		t.Fatal("expected parent key, got nil")
	}

	t.Run("panic_on_invalid_id", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic for invalid assetID")
			}
		}()
		NewAssetKey(spaceID, "invalid ID!")
	})
}

func TestNewAssetEntryWithoutID(t *testing.T) {
	spaceID := coretypes.SpaceID("testspace")
	entry := NewAssetEntryWithoutID(spaceID)
	if entry.Key == nil {
		t.Error("expected entry.Key, got nil")
	}
	if entry.Data == nil {
		t.Error("expected entry.Data, got nil")
	}
	if entry.Record == nil {
		t.Error("expected entry.Record, got nil")
	}
}

func TestNewAssetEntry(t *testing.T) {
	spaceID := coretypes.SpaceID("testspace")
	assetID := "test_asset"
	entry := NewAssetEntry(spaceID, assetID)
	if entry.ID != assetID {
		t.Errorf("expected ID %s, got %s", assetID, entry.ID)
	}
	if entry.FullID != string(spaceID)+":"+assetID {
		t.Errorf("expected FullID %s:%s, got %s", spaceID, assetID, entry.FullID)
	}
	if entry.Key == nil {
		t.Error("expected entry.Key, got nil")
	}
	if entry.Data == nil {
		t.Error("expected entry.Data, got nil")
	}
	if entry.Record == nil {
		t.Error("expected entry.Record, got nil")
	}
}
