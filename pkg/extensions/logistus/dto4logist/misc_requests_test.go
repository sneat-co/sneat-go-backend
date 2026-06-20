package dto4logist

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/with"
)

func TestCreateCounterpartyRequest_Validate(t *testing.T) {
	validCompany := dto4contactus.CreateCompanyRequest{Title: "Acme", CountryID: "US"}
	tests := []struct {
		name    string
		v       CreateCounterpartyRequest
		wantErr bool
	}{
		{"valid", CreateCounterpartyRequest{
			SpaceRequest: dto4spaceus.ValidSpaceRequest(),
			RolesField:   with.RolesField{Roles: []string{"buyer"}},
			Company:      validCompany,
		}, false},
		{"bad_space", CreateCounterpartyRequest{
			RolesField: with.RolesField{Roles: []string{"buyer"}},
			Company:    validCompany,
		}, true},
		{"bad_company", CreateCounterpartyRequest{
			SpaceRequest: dto4spaceus.ValidSpaceRequest(),
			RolesField:   with.RolesField{Roles: []string{"buyer"}},
			Company:      dto4contactus.CreateCompanyRequest{},
		}, true},
		{"duplicate_roles", CreateCounterpartyRequest{
			SpaceRequest: dto4spaceus.ValidSpaceRequest(),
			RolesField:   with.RolesField{Roles: []string{"buyer", "buyer"}},
			Company:      validCompany,
		}, true},
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

func TestNewOrderRequest(t *testing.T) {
	r := NewOrderRequest("space1", "order1")
	assert.Equal(t, "order1", r.OrderID)
	assert.NoError(t, r.Validate())
}

func TestOrderResponse_Validate(t *testing.T) {
	assert.Error(t, OrderResponse{}.Validate())
	assert.NoError(t, OrderResponse{OrderDto: &dbo4logist.OrderDbo{}}.Validate())
}

func TestValidateID(t *testing.T) {
	assert.NoError(t, validateID("f", "abc"))
	assert.Error(t, validateID("f", ""))
	assert.Error(t, validateID("f", "  "))
	assert.Error(t, validateID("f", "a b"))
	assert.Error(t, validateID("f", "a\tb"))
}

func TestValidateContainerID(t *testing.T) {
	assert.NoError(t, validateContainerID("containerID", "c1"))
	assert.Error(t, validateContainerID("containerID", ""))
}

func TestDeleteOrderCounterpartyRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       DeleteOrderCounterpartyRequest
		wantErr bool
	}{
		{"valid", DeleteOrderCounterpartyRequest{
			OrderRequest: ValidOrderRequest(),
			Role:         dbo4logist.CounterpartyRoleBuyer,
			ContactID:    "c1",
		}, false},
		{"bad_order", DeleteOrderCounterpartyRequest{Role: dbo4logist.CounterpartyRoleBuyer, ContactID: "c1"}, true},
		{"missing_role", DeleteOrderCounterpartyRequest{OrderRequest: ValidOrderRequest(), ContactID: "c1"}, true},
		{"missing_contact", DeleteOrderCounterpartyRequest{OrderRequest: ValidOrderRequest(), Role: dbo4logist.CounterpartyRoleBuyer}, true},
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

func TestSetOrderStatusRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       SetOrderStatusRequest
		wantErr bool
	}{
		{"valid", SetOrderStatusRequest{OrderRequest: ValidOrderRequest(), Status: dbo4logist.OrderStatusDraft}, false},
		{"bad_space", SetOrderStatusRequest{Status: dbo4logist.OrderStatusDraft}, true},
		{"missing_status", SetOrderStatusRequest{OrderRequest: ValidOrderRequest()}, true},
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

func TestSetOrderCounterparty_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       SetOrderCounterparty
		wantErr bool
	}{
		{"valid", SetOrderCounterparty{ContactID: "c1", Role: "buyer"}, false},
		{"missing_contact", SetOrderCounterparty{Role: "buyer"}, true},
		{"missing_role", SetOrderCounterparty{ContactID: "c1"}, true},
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

func TestSetOrderCounterpartiesRequest_Validate(t *testing.T) {
	valid := SetOrderCounterparty{ContactID: "c1", Role: "buyer"}
	tests := []struct {
		name    string
		v       SetOrderCounterpartiesRequest
		wantErr bool
	}{
		{"valid", SetOrderCounterpartiesRequest{
			OrderRequest:   ValidOrderRequest(),
			Counterparties: []SetOrderCounterparty{valid},
		}, false},
		{"bad_order", SetOrderCounterpartiesRequest{Counterparties: []SetOrderCounterparty{valid}}, true},
		{"no_counterparties", SetOrderCounterpartiesRequest{OrderRequest: ValidOrderRequest()}, true},
		{"bad_counterparty", SetOrderCounterpartiesRequest{
			OrderRequest:   ValidOrderRequest(),
			Counterparties: []SetOrderCounterparty{{}},
		}, true},
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

func TestOrderShippingPointRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       OrderShippingPointRequest
		wantErr bool
	}{
		{"valid", OrderShippingPointRequest{OrderRequest: ValidOrderRequest(), ShippingPointID: "sp1"}, false},
		{"bad_order", OrderShippingPointRequest{ShippingPointID: "sp1"}, true},
		{"bad_shipping_point", OrderShippingPointRequest{OrderRequest: ValidOrderRequest(), ShippingPointID: " sp1 "}, true},
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

func TestUpdateShippingPointRequest_Validate(t *testing.T) {
	validBase := OrderShippingPointRequest{OrderRequest: ValidOrderRequest(), ShippingPointID: "sp1"}
	tests := []struct {
		name    string
		v       UpdateShippingPointRequest
		wantErr bool
	}{
		{"valid_string", UpdateShippingPointRequest{
			OrderShippingPointRequest: validBase,
			SetFieldsRequest:          SetFieldsRequest{SetStrings: map[string]string{"notes": "hi"}},
		}, false},
		{"valid_date", UpdateShippingPointRequest{
			OrderShippingPointRequest: validBase,
			SetFieldsRequest:          SetFieldsRequest{SetDates: map[string]string{"scheduledStartDate": "2023-01-01"}},
		}, false},
		{"bad_base", UpdateShippingPointRequest{
			SetFieldsRequest: SetFieldsRequest{SetStrings: map[string]string{"notes": "hi"}},
		}, true},
		{"nothing_set", UpdateShippingPointRequest{OrderShippingPointRequest: validBase}, true},
		{"unknown_date", UpdateShippingPointRequest{
			OrderShippingPointRequest: validBase,
			SetFieldsRequest:          SetFieldsRequest{SetDates: map[string]string{"bogus": "2023-01-01"}},
		}, true},
		{"unknown_string", UpdateShippingPointRequest{
			OrderShippingPointRequest: validBase,
			SetFieldsRequest:          SetFieldsRequest{SetStrings: map[string]string{"bogus": "x"}},
		}, true},
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

func TestAddContainerPoint_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       AddContainerPoint
		wantErr bool
	}{
		{"valid", AddContainerPoint{ID: "c1", Tasks: []task{dbo4logist.ShippingPointTaskLoad}}, false},
		{"missing_id", AddContainerPoint{Tasks: []task{dbo4logist.ShippingPointTaskLoad}}, true},
		{"no_tasks", AddContainerPoint{ID: "c1"}, true},
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

func TestAddOrderShippingPointRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       AddOrderShippingPointRequest
		wantErr bool
	}{
		{"valid_with_container", AddOrderShippingPointRequest{
			OrderRequest:      ValidOrderRequest(),
			LocationContactID: "loc1",
			Containers:        []AddContainerPoint{{ID: "c1", Tasks: []task{dbo4logist.ShippingPointTaskLoad}}},
		}, false},
		{"bad_order", AddOrderShippingPointRequest{
			LocationContactID: "loc1",
			Tasks:             []task{dbo4logist.ShippingPointTaskLoad},
		}, true},
		{"missing_location", AddOrderShippingPointRequest{
			OrderRequest: ValidOrderRequest(),
			Tasks:        []task{dbo4logist.ShippingPointTaskLoad},
		}, true},
		{"bad_container", AddOrderShippingPointRequest{
			OrderRequest:      ValidOrderRequest(),
			LocationContactID: "loc1",
			Containers:        []AddContainerPoint{{ID: "c1"}},
		}, true},
		{"duplicate_container", AddOrderShippingPointRequest{
			OrderRequest:      ValidOrderRequest(),
			LocationContactID: "loc1",
			Containers: []AddContainerPoint{
				{ID: "c1", Tasks: []task{dbo4logist.ShippingPointTaskLoad}},
				{ID: "c1", Tasks: []task{dbo4logist.ShippingPointTaskLoad}},
			},
		}, true},
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

func TestSetLogistSpaceSettingsRequest_Validate(t *testing.T) {
	validAddr := dbmodels.Address{CountryID: "US"}
	tests := []struct {
		name    string
		v       SetLogistSpaceSettingsRequest
		wantErr bool
	}{
		{"valid", SetLogistSpaceSettingsRequest{
			SpaceRequest: dto4spaceus.ValidSpaceRequest(),
			Roles:        []dbo4logist.LogistSpaceRole{dbo4logist.CompanyRoleTrucker},
			Address:      validAddr,
		}, false},
		{"bad_space", SetLogistSpaceSettingsRequest{
			Roles:   []dbo4logist.LogistSpaceRole{dbo4logist.CompanyRoleTrucker},
			Address: validAddr,
		}, true},
		{"no_roles", SetLogistSpaceSettingsRequest{
			SpaceRequest: dto4spaceus.ValidSpaceRequest(),
			Address:      validAddr,
		}, true},
		{"unknown_role", SetLogistSpaceSettingsRequest{
			SpaceRequest: dto4spaceus.ValidSpaceRequest(),
			Roles:        []dbo4logist.LogistSpaceRole{"bogus-role"},
			Address:      validAddr,
		}, true},
		{"bad_address", SetLogistSpaceSettingsRequest{
			SpaceRequest: dto4spaceus.ValidSpaceRequest(),
			Roles:        []dbo4logist.LogistSpaceRole{dbo4logist.CompanyRoleTrucker},
			Address:      dbmodels.Address{},
		}, true},
		{"untrimmed_vat", SetLogistSpaceSettingsRequest{
			SpaceRequest: dto4spaceus.ValidSpaceRequest(),
			Roles:        []dbo4logist.LogistSpaceRole{dbo4logist.CompanyRoleTrucker},
			Address:      validAddr,
			VATNumber:    " VAT ",
		}, true},
		{"vat_too_long", SetLogistSpaceSettingsRequest{
			SpaceRequest: dto4spaceus.ValidSpaceRequest(),
			Roles:        []dbo4logist.LogistSpaceRole{dbo4logist.CompanyRoleTrucker},
			Address:      validAddr,
			VATNumber:    "123456789012345678901",
		}, true},
		{"untrimmed_prefix", SetLogistSpaceSettingsRequest{
			SpaceRequest:      dto4spaceus.ValidSpaceRequest(),
			Roles:             []dbo4logist.LogistSpaceRole{dbo4logist.CompanyRoleTrucker},
			Address:           validAddr,
			OrderNumberPrefix: " AB ",
		}, true},
		{"prefix_too_long", SetLogistSpaceSettingsRequest{
			SpaceRequest:      dto4spaceus.ValidSpaceRequest(),
			Roles:             []dbo4logist.LogistSpaceRole{dbo4logist.CompanyRoleTrucker},
			Address:           validAddr,
			OrderNumberPrefix: "ABCDEF",
		}, true},
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
