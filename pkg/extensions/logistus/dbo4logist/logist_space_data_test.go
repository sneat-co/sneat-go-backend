package dbo4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func TestOrderCounter_Validate(t *testing.T) {
	assert.NoError(t, OrderCounter{Prefix: "X", LastNumber: 0}.Validate())
	assert.NoError(t, OrderCounter{LastNumber: 5}.Validate())
	assert.Error(t, OrderCounter{LastNumber: -1}.Validate())
}

func TestLogistSpaceDbo_Validate(t *testing.T) {
	withUsers := dbmodels.WithUserIDs{UserIDs: []string{"u1"}}
	tests := []struct {
		name    string
		v       LogistSpaceDbo
		wantErr bool
	}{
		{"valid", LogistSpaceDbo{WithUserIDs: withUsers, Roles: []string{string(CompanyRoleBuyer)}}, false},
		{"missing_users", LogistSpaceDbo{Roles: []string{string(CompanyRoleBuyer)}}, true},
		{"missing_roles", LogistSpaceDbo{WithUserIDs: withUsers}, true},
		{"unknown_role", LogistSpaceDbo{WithUserIDs: withUsers, Roles: []string{"bad_role"}}, true},
		{"contact_untrimmed", LogistSpaceDbo{WithUserIDs: withUsers, Roles: []string{string(CompanyRoleBuyer)}, ContactID: " c "}, true},
		{"prefix_untrimmed", LogistSpaceDbo{WithUserIDs: withUsers, Roles: []string{string(CompanyRoleBuyer)}, OrderNumberPrefix: " P "}, true},
		{"prefix_too_long", LogistSpaceDbo{WithUserIDs: withUsers, Roles: []string{string(CompanyRoleBuyer)}, OrderNumberPrefix: "ABCDEF"}, true},
		{"bad_counter", LogistSpaceDbo{WithUserIDs: withUsers, Roles: []string{string(CompanyRoleBuyer)},
			OrderCounters: map[string]OrderCounter{"k": {LastNumber: -1}}}, true},
		{"valid_counter", LogistSpaceDbo{WithUserIDs: withUsers, Roles: []string{string(CompanyRoleBuyer)},
			OrderCounters: map[string]OrderCounter{"k": {LastNumber: 1}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
