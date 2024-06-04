package dbo4logist

import (
	"fmt"
	"github.com/strongo/validation"
)

// OrderRoute is a route for an order
type OrderRoute struct {
	Origin        TransitPoint    `json:"origin" firestore:"origin"`
	Destination   TransitPoint    `json:"destination" firestore:"destination"`
	TransitPoints []*TransitPoint `json:"transitPoints,omitempty" firestore:"transitPoints,omitempty"`
}

// Validate returns error if order route is invalid
func (v OrderRoute) Validate() error {
	if err := v.Origin.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("origin", err.Error())
	}
	if err := v.Destination.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("destination", err.Error())
	}
	for i, tp := range v.TransitPoints {
		if err := tp.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("transit[%v]", i), err.Error())
		}
	}
	return nil
}
