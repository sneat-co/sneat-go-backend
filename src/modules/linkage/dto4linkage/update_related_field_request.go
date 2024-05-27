package dto4linkage

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/strongo/validation"
)

// UpdateRelatedFieldRequest is a request to update related items
type UpdateRelatedFieldRequest struct {
	Related map[string]*models4linkage.RelationshipRolesCommand `json:"related"`
}

// Validate checks if request is valid
func (v *UpdateRelatedFieldRequest) Validate() error {
	for id, rel := range v.Related {
		if err := rel.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("related."+id, err.Error())
		}
	}
	return nil
}
