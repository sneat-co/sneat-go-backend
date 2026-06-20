package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertLogistSpaceRolesToStringSlice(t *testing.T) {
	roles := []LogistSpaceRole{CompanyRoleBuyer, CompanyRoleTrucker}
	got := ConvertLogistSpaceRolesToStringSlice(roles)
	assert.Equal(t, []string{string(CompanyRoleBuyer), string(CompanyRoleTrucker)}, got)
	assert.Empty(t, ConvertLogistSpaceRolesToStringSlice(nil))
}

func TestRolesChanged(t *testing.T) {
	current := []string{string(CompanyRoleBuyer), string(CompanyRoleTrucker)}
	assert.False(t, RolesChanged(current, []LogistSpaceRole{CompanyRoleBuyer}))
	assert.True(t, RolesChanged(current, []LogistSpaceRole{CompanyRoleShippingLine}))
	assert.False(t, RolesChanged(current, nil))
}

func TestIsKnownLogistCompanyRole(t *testing.T) {
	assert.True(t, IsKnownLogistCompanyRole(CompanyRoleBuyer))
	assert.True(t, IsKnownLogistCompanyRole(CompanyRoleWarehouseOperator))
	assert.False(t, IsKnownLogistCompanyRole("not_a_role"))
}
