package dbo4assetus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/briefs4assetus"
	"github.com/strongo/validation"
)

type AssetusSpaceBrief struct { // TODO: document intended usage & provide example use cases
	briefs4assetus.WithAssets
}

type WithAssetSpaces struct { // TODO: document intended usage & provide example use cases
	Spaces map[string]*AssetusSpaceBrief `json:"spaces,omitempty" firestore:"spaces,omitempty"`
}

// Validate returns error if not valid
func (v *WithAssetSpaces) Validate() error {
	for id, assetusSpaceBrief := range v.Spaces {
		if id == "" {
			return validation.NewErrBadRecordFieldValue("spaces", "spaceID can not be empty string")
		}
		if assetusSpaceBrief == nil {
			return validation.NewErrBadRecordFieldValue("spaces."+id, "can not be nil")
		}
		if err := assetusSpaceBrief.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("spaces."+id, err.Error())
		}
	}
	return nil
}

// AddAsset adds an asset to a team
func (v *WithAssetSpaces) AddAsset(teamID string, asset AssetEntry) (updates []dal.Update, err error) {
	if v.Spaces == nil {
		v.Spaces = make(map[string]*AssetusSpaceBrief)
	}
	assetusSpaceBrief := v.Spaces[teamID]
	if assetusSpaceBrief == nil {
		assetusSpaceBrief = new(AssetusSpaceBrief)
		v.Spaces[teamID] = assetusSpaceBrief
	}
	if assetusSpaceBrief.Assets == nil {
		assetusSpaceBrief.Assets = make(map[string]*briefs4assetus.AssetBrief)
	}
	var assetBrief briefs4assetus.AssetBrief
	if assetBrief, err = asset.Data.GetAssetBrief(); err != nil {
		return
	}
	if updates, err = assetusSpaceBrief.AddAssetBrief(asset.ID, assetBrief); err != nil {
		return
	}
	for i, u := range updates {
		u.Field = "spaces." + teamID + "." + u.Field
		updates[i] = u
	}
	return
}
