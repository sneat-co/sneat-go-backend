package dbo4assetus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"
	"github.com/strongo/validation"
)

type AssetusTeamBrief struct { // TODO: document intended usage & provide example use cases
	briefs4assetus.WithAssets
}

type WithAssetTeams struct { // TODO: document intended usage & provide example use cases
	Teams map[string]*AssetusTeamBrief `json:"teams,omitempty" firestore:"teams,omitempty"`
}

// Validate returns error if not valid
func (v *WithAssetTeams) Validate() error {
	for id, assetusTeamBrief := range v.Teams {
		if id == "" {
			return validation.NewErrBadRecordFieldValue("teams", "teamID can not be empty string")
		}
		if assetusTeamBrief == nil {
			return validation.NewErrBadRecordFieldValue("teams."+id, "can not be nil")
		}
		if err := assetusTeamBrief.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("teams."+id, err.Error())
		}
	}
	return nil
}

// AddAsset adds an asset to a team
func (v *WithAssetTeams) AddAsset(teamID string, asset AssetEntry) (updates []dal.Update, err error) {
	if v.Teams == nil {
		v.Teams = make(map[string]*AssetusTeamBrief)
	}
	assetusTeamBrief := v.Teams[teamID]
	if assetusTeamBrief == nil {
		assetusTeamBrief = new(AssetusTeamBrief)
		v.Teams[teamID] = assetusTeamBrief
	}
	if assetusTeamBrief.Assets == nil {
		assetusTeamBrief.Assets = make(map[string]*briefs4assetus.AssetBrief)
	}
	var assetBrief briefs4assetus.AssetBrief
	if assetBrief, err = asset.Data.GetAssetBrief(); err != nil {
		return
	}
	if updates, err = assetusTeamBrief.AddAssetBrief(asset.ID, assetBrief); err != nil {
		return
	}
	for i, u := range updates {
		u.Field = "teams." + teamID + "." + u.Field
		updates[i] = u
	}
	return
}
