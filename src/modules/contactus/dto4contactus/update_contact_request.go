package dto4contactus

import (
	"errors"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

type UpdateRelatedRequest struct {
	Related map[string]*models4linkage.RelationshipRolesCommand `json:"related"`
}

func (v *UpdateRelatedRequest) Validate() error {
	for id, rel := range v.Related {
		if err := rel.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("related."+id, err.Error())
		}
	}
	return nil
}

type UpdateContactRequest struct {
	ContactRequest
	UpdateRelatedRequest
	Address   *dbmodels.Address `json:"address,omitempty"`
	AgeGroup  string            `json:"ageGroup,omitempty"`
	Roles     *SetRolesRequest  `json:"roles,omitempty"`
	VatNumber *string           `json:"vatNumber,omitempty"`
}

func (v UpdateContactRequest) Validate() error {
	if err := v.ContactRequest.Validate(); err != nil {
		return err
	}
	if v.Address == nil && v.AgeGroup == "" && v.Roles == nil && v.Related == nil && v.VatNumber == nil {
		return validation.NewBadRequestError(errors.New("at least one of contact fields must be provided for an update"))
	}
	if v.Address != nil {
		if err := v.Address.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("address", err.Error())
		}
	}
	if err := dbmodels.ValidateAgeGroup(v.AgeGroup, false); err != nil {
		return validation.NewErrBadRequestFieldValue("ageGroup", err.Error())
	}
	if v.Roles != nil {
		if err := v.Roles.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("roles", err.Error())
		}
	}
	if v.VatNumber != nil {
		vat := *v.VatNumber
		if strings.TrimSpace(vat) == vat {
			return validation.NewErrBadRequestFieldValue("vatNumber", "must not have leading or trailing spaces")
		}

	}
	if v.Related != nil {
		if err := v.UpdateRelatedRequest.Validate(); err != nil {
			return err
		}
	}
	return nil
}
