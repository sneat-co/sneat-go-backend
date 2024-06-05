package extra

var factories = make(map[Type]func() Data)

func newExtra(extraType Type) Data {
	if factory, ok := factories[extraType]; ok {
		return factory()
	}
	return new(noExtra)
}

func RegisterFactory(extraType Type, factory func() Data) {
	factories[extraType] = factory
}
