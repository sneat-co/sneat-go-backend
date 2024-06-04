package extras4assetus

var assetExtraFactories = map[AssetExtraType]func() AssetExtra{}

func RegisterAssetExtraFactory(t AssetExtraType, f func() AssetExtra) {
	assetExtraFactories[t] = f
}

func NewAssetExtra(t AssetExtraType) AssetExtra {
	if f, ok := assetExtraFactories[t]; ok {
		return f()
	}
	return nil
}
