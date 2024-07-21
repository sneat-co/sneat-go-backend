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

// WithAssets // TODO: should be moved to assetus module?
type WithAssets struct {
	Assets map[string]*AssetBrief `json:"assets,omitempty" firestore:"assets,omitempty"`
}

func (v *WithAssets) AddAssetBrief(assetID string, assetBrief AssetBrief) (updates []dal.Update, err error) {
	if v.Assets == nil {
		v.Assets = make(map[string]*AssetBrief, 1)
	}
	v.Assets[assetID] = &assetBrief
	updates = append(updates, dal.Update{Field: "assets." + assetID, Value: assetBrief})
	return
}

func (v *WithAssets) GetAssetBriefByID(id string) (brief *AssetBrief) {
	return v.Assets[id]
}

// Validate returns error if not valid
func (v *WithAssets) Validate() error {
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
