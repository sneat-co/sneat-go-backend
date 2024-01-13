package models4linkage

import core "github.com/sneat-co/sneat-go-core"

type Relatable interface {
	GetRelated() *WithRelatedAndIDs
	core.Validatable
}
