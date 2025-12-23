package dbo4logist

import (
	"slices"

	"github.com/strongo/slice"
)

type LogistSpaceRole string

func ConvertLogistSpaceRolesToStringSlice(roles []LogistSpaceRole) []string {
	result := make([]string, len(roles))
	for i, r := range roles {
		result[i] = string(r)
	}
	return result
}

func RolesChanged(currentRoles []string, newRoles []LogistSpaceRole) bool {
	for _, r := range newRoles {
		if role := string(r); !slices.Contains(currentRoles, role) {
			return true
		}
	}
	return false
}

const (
	CompanyRoleBuyer        = LogistSpaceRole(CounterpartyRoleBuyer)
	CompanyRoleCustomBroker = LogistSpaceRole(CounterpartyRoleCustomBroker)
	CompanyRoleDispatcher   = LogistSpaceRole(CounterpartyRoleDispatcher)
	CompanyRoleTrucker      = LogistSpaceRole(CounterpartyRoleTrucker)
	CompanyRoleShippingLine = LogistSpaceRole(CounterpartyRoleShippingLine)

	CompanyRoleFreightAgent  = "freight_agent" // Freight agent can be either a CounterpartyRoleDispatchAgent or CounterpartyRoleReceiveAgent
	CompanyRoleDispatchAgent = LogistSpaceRole(CounterpartyRoleDispatchAgent)
	CompanyRoleReceiveAgent  = LogistSpaceRole(CounterpartyRoleReceiveAgent)

	CompanyRoleFreightBroker     LogistSpaceRole = "freight_broker"
	CompanyRoleFreightForwarder  LogistSpaceRole = "freight_forwarder"
	CompanyRoleWarehouseOperator LogistSpaceRole = "warehouse_operator"
)

var KnownLogistCompanyRoles = []LogistSpaceRole{
	CompanyRoleBuyer,
	CompanyRoleCustomBroker,
	CompanyRoleFreightAgent,
	CompanyRoleFreightBroker,
	CompanyRoleFreightForwarder,
	CompanyRoleDispatchAgent,
	CompanyRoleReceiveAgent,
	CompanyRoleDispatcher,
	CompanyRoleShippingLine,
	CompanyRoleTrucker,
	CompanyRoleWarehouseOperator,
}

func IsKnownLogistCompanyRole(role LogistSpaceRole) bool {
	return slice.Index(KnownLogistCompanyRoles, role) >= 0
}
