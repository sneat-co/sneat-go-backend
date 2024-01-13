package briefs4assetus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// Asset is an interface restriction for asset briefs to use in generic types
type Asset interface {
	dbmodels.RelatedAs
}

type AssetusTeamBrief[T Asset] struct {
	WithAssets[T]
}

type WithAssetusTeamBriefs[T Asset] struct {
	Teams map[string]*AssetusTeamBrief[T] `json:"teams,omitempty" firestore:"teams,omitempty"`
}

// WithAssets // TODO: should be moved to assetus module?
type WithAssets[T Asset] struct {
	Assets map[string]T `json:"assets,omitempty" firestore:"assets,omitempty"`
}

func (v *WithAssets[T]) AddAsset(assetID string, assetBrief T) (updates []dal.Update) {
	if v.Assets == nil {
		v.Assets = make(map[string]T, 1)
	}
	v.Assets[assetID] = assetBrief
	updates = append(updates, dal.Update{Field: "assets." + assetID, Value: assetBrief})
	return
}

func (v *WithAssets[T]) GetAssetBriefByID(id string) (brief T) {
	return v.Assets[id]
}

// Validate returns error if not valid
func (v *WithAssets[T]) Validate() error {
	for id, asset := range v.Assets {
		if id == "" {
			return validation.NewErrBadRecordFieldValue("assets", "assetID can not be empty string")
		}
		if err := asset.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("assets."+id, err.Error())
		}
	}
	return nil
}

// Validate returns error if not valid
func (v *WithAssetusTeamBriefs[T]) Validate() error {
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
func (v *WithAssetusTeamBriefs[T]) AddAsset(teamID, assetID string, assetBrief T) (updates []dal.Update) {
	if v.Teams == nil {
		v.Teams = make(map[string]*AssetusTeamBrief[T])
	}
	assetusTeamBrief := v.Teams[teamID]
	if assetusTeamBrief == nil {
		assetusTeamBrief = new(AssetusTeamBrief[T])
		v.Teams[teamID] = assetusTeamBrief
	}
	if assetusTeamBrief.Assets == nil {
		assetusTeamBrief.Assets = make(map[string]T)
	}
	updates = assetusTeamBrief.AddAsset(assetID, assetBrief)
	for i, u := range updates {
		u.Field = "teams." + teamID + "." + u.Field
		updates[i] = u
	}
	return updates
}
