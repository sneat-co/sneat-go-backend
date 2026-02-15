package extras4assetus

import (
	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/briefs4assetus"
)

var assetExtraFactories = map[extra.Type]func() briefs4assetus.AssetExtra{}

func RegisterAssetExtraFactory(t extra.Type, f func() briefs4assetus.AssetExtra) {
	assetExtraFactories[t] = f
}

func NewAssetExtra(t extra.Type) briefs4assetus.AssetExtra {
	if f, ok := assetExtraFactories[t]; ok {
		return f()
	}
	return nil
}
