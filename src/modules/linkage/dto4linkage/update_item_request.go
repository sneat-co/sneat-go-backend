package dto4linkage

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
)

type UpdateItemRequest struct {
	dbo4linkage.TeamModuleItemRef `json:"itemRef"`
	UpdateRelatedFieldRequest
}

func (v *UpdateItemRequest) Validate() error {
	if err := v.TeamModuleItemRef.Validate(); err != nil {
		return err
	}
	if err := v.UpdateRelatedFieldRequest.Validate(); err != nil {
		return err
	}
	return nil
}
