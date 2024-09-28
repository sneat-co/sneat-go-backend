package entities

import "fmt"

var registered = map[string]Entity{}

// Register entity
func Register(entity Entity) {
	if v, ok := registered[entity.Name]; ok {
		panic(fmt.Sprintf("duplicate entity name=%s. Already registered: %T; Registering: %T;",
			entity.Name, v, entity))
	}
	registered[entity.Name] = entity
}
