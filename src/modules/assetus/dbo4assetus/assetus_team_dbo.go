package dbo4assetus

import "github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"

// AssetusTeamDbo summary about assets for a team
type AssetusTeamDbo struct {
	briefs4assetus.WithAssets
	WithAssetTeams
}

// Validate returns error if not valid
func (v AssetusTeamDbo) Validate() error {
	if err := v.WithAssets.Validate(); err != nil {
		return err
	}
	if err := v.WithAssetTeams.Validate(); err != nil {
		return err
	}
	return nil
}
