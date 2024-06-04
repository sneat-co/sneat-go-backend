package dto4logist

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

// SetLogistTeamSettingsRequest represents a request to set logistus team settings
type SetLogistTeamSettingsRequest struct {
	dto4teamus.TeamRequest
	Roles             []dbo4logist.LogistTeamRole `json:"roles"`
	Address           dbmodels.Address            `json:"address"`
	VATNumber         string                      `json:"vatNumber,omitempty"`
	OrderNumberPrefix string                      `json:"orderNumberPrefix,omitempty"`
}

// Validate returns error if request is invalid
func (v SetLogistTeamSettingsRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if len(v.Roles) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("roles")
	}
	for i, role := range v.Roles {
		if !dbo4logist.IsKnownLogistCompanyRole(role) {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("roles[%d]", i),
				fmt.Sprintf("should be one of: %+v", dbo4logist.KnownLogistCompanyRoles))
		}
	}
	if err := v.Address.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.VATNumber) != v.VATNumber {
		return validation.NewErrBadRequestFieldValue("address", "should be trimmed")
	}
	if len(v.VATNumber) > 20 {
		return validation.NewErrBadRequestFieldValue("vatNumber", "should not be longer than 20 characters")
	}
	if strings.TrimSpace(v.OrderNumberPrefix) != v.OrderNumberPrefix {
		return validation.NewErrBadRequestFieldValue("orderNumberPrefix", "should be trimmed")
	}
	if len(v.OrderNumberPrefix) > 5 {
		return validation.NewErrBadRequestFieldValue("vatNumber", "should not be longer than 5 characters")
	}
	return nil
}
