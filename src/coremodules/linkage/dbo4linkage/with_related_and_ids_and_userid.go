package dbo4linkage

import (
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

type WithRelatedAndIDsAndUserID struct {
	dbmodels.WithUserID
	*WithRelatedAndIDs
}

func (v *WithRelatedAndIDsAndUserID) Validate() error {
	if err := v.WithUserID.Validate(); err != nil {
		return err
	}
	if err := v.WithRelatedAndIDs.Validate(); err != nil {
		return err
	}
	return nil
}
