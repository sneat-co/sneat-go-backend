package extras4assetus

import "github.com/sneat-co/sneat-go-backend/src/coremodels/extra"

var assetExtraFactories = map[extra.Type]func() extra.Data{}

func RegisterAssetExtraFactory(t extra.Type, f func() extra.Data) {
	assetExtraFactories[t] = f
}

func NewAssetExtra(t extra.Type) extra.Data {
	if f, ok := assetExtraFactories[t]; ok {
		return f()
	}
	return nil
}
