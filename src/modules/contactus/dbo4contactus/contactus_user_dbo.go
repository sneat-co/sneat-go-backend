package models4contactus

import "github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"

// ContactusUserDbo holds contactus specific data for a specific user
type ContactusUserDbo struct {
	dbo4linkage.WithRelated
}

// Validate returns error if not valid
func (v *ContactusUserDbo) Validate() error {
	if err := v.WithRelated.Validate(); err != nil {
		return err
	}
	return nil
}
