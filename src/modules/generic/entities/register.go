package entities

import "fmt"

var registered = map[string]Entity{}

// Register entity
func Register(entity Entity) {
	if v, ok := registered[entity.Name]; ok {
		panic(fmt.Sprintf("duplicate entity name=%v. Already registered: %T; Registering: %T;",
			entity.Name, v, entity))
	}
	registered[entity.Name] = entity
}
