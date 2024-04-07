package models4logist

import (
	"github.com/strongo/validation"
	"strings"
)

// TransitPoint is a transit point of an order route
type TransitPoint struct {
	CountryID string `json:"countryID" firestore:"countryID"`
	Address   string `json:"address,omitempty" firestore:"address,omitempty"`
}

// Validate returns error if transit point is invalid
func (v TransitPoint) Validate() error {
	if strings.TrimSpace(v.CountryID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("countryID")
	}
	return nil
}
