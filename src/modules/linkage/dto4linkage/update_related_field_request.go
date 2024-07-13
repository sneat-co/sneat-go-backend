package dto4linkage

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/strongo/validation"
)

// UpdateRelatedFieldRequest is a request to update related items
type UpdateRelatedFieldRequest struct {
	Related []dbo4linkage.RelationshipItemRolesCommand `json:"related"`
}

// Validate checks if request is valid
func (v *UpdateRelatedFieldRequest) Validate() error {
	for i, rel := range v.Related {
		if err := rel.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("related[%d].", i), err.Error())
		}
	}
	return nil
}
