package dbo4assetus

import "github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"

// AssetusSpaceDbo summary about assets for a team
type AssetusSpaceDbo struct {
	briefs4assetus.WithAssets
	WithAssetSpaces
}

// Validate returns error if not valid
func (v AssetusSpaceDbo) Validate() error {
	if err := v.WithAssets.Validate(); err != nil {
		return err
	}
	if err := v.WithAssetSpaces.Validate(); err != nil {
		return err
	}
	return nil
}
