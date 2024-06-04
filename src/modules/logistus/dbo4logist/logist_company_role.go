package dbo4logist

import "github.com/strongo/slice"

type LogistTeamRole string

func ConvertLogistTeamRolesToStringSlice(roles []LogistTeamRole) []string {
	result := make([]string, len(roles))
	for i, r := range roles {
		result[i] = string(r)
	}
	return result
}

func RolesChanged(currentRoles []string, newRoles []LogistTeamRole) bool {
	for _, r := range newRoles {
		if role := string(r); !slice.Contains(currentRoles, role) {
			return true
		}
	}
	return false
}

const (
	CompanyRoleBuyer        = LogistTeamRole(CounterpartyRoleBuyer)
	CompanyRoleCustomBroker = LogistTeamRole(CounterpartyRoleCustomBroker)
	CompanyRoleDispatcher   = LogistTeamRole(CounterpartyRoleDispatcher)
	CompanyRoleTrucker      = LogistTeamRole(CounterpartyRoleTrucker)
	CompanyRoleShippingLine = LogistTeamRole(CounterpartyRoleShippingLine)

	CompanyRoleFreightAgent  = "freight_agent" // Freight agent can be either a CounterpartyRoleDispatchAgent or CounterpartyRoleReceiveAgent
	CompanyRoleDispatchAgent = LogistTeamRole(CounterpartyRoleDispatchAgent)
	CompanyRoleReceiveAgent  = LogistTeamRole(CounterpartyRoleReceiveAgent)

	CompanyRoleFreightBroker     LogistTeamRole = "freight_broker"
	CompanyRoleFreightForwarder  LogistTeamRole = "freight_forwarder"
	CompanyRoleWarehouseOperator LogistTeamRole = "warehouse_operator"
)

var KnownLogistCompanyRoles = []LogistTeamRole{
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

func IsKnownLogistCompanyRole(role LogistTeamRole) bool {
	return slice.Index(KnownLogistCompanyRoles, role) >= 0
}
