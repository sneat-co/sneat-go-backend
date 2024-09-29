package dto4linkage

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/linkage/dbo4linkage"
)

type UpdateItemRequest struct {
	dbo4linkage.SpaceModuleItemRef `json:"itemRef"`
	UpdateRelatedFieldRequest
}

func (v *UpdateItemRequest) Validate() error {
	if err := v.SpaceModuleItemRef.Validate(); err != nil {
		return err
	}
	if err := v.UpdateRelatedFieldRequest.Validate(); err != nil {
		return err
	}
	return nil
}
