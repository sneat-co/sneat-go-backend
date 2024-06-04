package dbo4assetus

import "github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"

// AssetusTeamDto summary about assets for a team
type AssetusTeamDto struct {
	briefs4assetus.WithAssets[*briefs4assetus.AssetBrief]
	briefs4assetus.WithAssetusTeamBriefs[*briefs4assetus.AssetBrief]
}

// Validate returns error if not valid
func (v AssetusTeamDto) Validate() error {
	if err := v.WithAssets.Validate(); err != nil {
		return err
	}
	if err := v.WithAssetusTeamBriefs.Validate(); err != nil {
		return err
	}
	return nil
}
