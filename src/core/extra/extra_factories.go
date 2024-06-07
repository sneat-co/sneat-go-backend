package extra

var factories = make(map[Type]func() Data)

// NewExtraData creates a new extra data object of the specified type
func NewExtraData(extraType Type) (extraData Data) {
	if factory, ok := factories[extraType]; ok {
		return factory()
	}
	// Default to no extra data if no factory is registered for the type
	//return &noExtra{BaseData{ExtraType: string(extraType)}}
	return new(noExtra)
}

func RegisterFactory(extraType Type, factory func() Data) {
	factories[extraType] = factory
}
