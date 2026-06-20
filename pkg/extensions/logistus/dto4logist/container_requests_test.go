package dto4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/stretchr/testify/assert"
)

func TestPointOfNewContainer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       PointOfNewContainer
		wantErr bool
	}{
		{"valid", PointOfNewContainer{ShippingPointID: "sp1"}, false},
		{"empty", PointOfNewContainer{}, true},
		{"blank", PointOfNewContainer{ShippingPointID: "   "}, true},
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

func validNewContainer() NewContainer {
	return NewContainer{
		OrderContainerBase: dbo4logist.OrderContainerBase{Type: dbo4logist.ContainerType20ft},
		Points:             []PointOfNewContainer{{ShippingPointID: "sp1"}},
	}
}

func TestNewContainer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       NewContainer
		wantErr bool
	}{
		{"valid", validNewContainer(), false},
		{"no_points_ok", NewContainer{OrderContainerBase: dbo4logist.OrderContainerBase{Type: dbo4logist.ContainerType20ft}}, false},
		{"bad_base", NewContainer{Points: []PointOfNewContainer{{ShippingPointID: "sp1"}}}, true},
		{"bad_point", NewContainer{
			OrderContainerBase: dbo4logist.OrderContainerBase{Type: dbo4logist.ContainerType20ft},
			Points:             []PointOfNewContainer{{}},
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

func TestAddContainersRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       AddContainersRequest
		wantErr bool
	}{
		{"valid", AddContainersRequest{
			OrderRequest: ValidOrderRequest(),
			Containers:   []NewContainer{validNewContainer()},
		}, false},
		{"bad_order", AddContainersRequest{
			Containers: []NewContainer{validNewContainer()},
		}, true},
		{"no_containers", AddContainersRequest{
			OrderRequest: ValidOrderRequest(),
		}, true},
		{"bad_container", AddContainersRequest{
			OrderRequest: ValidOrderRequest(),
			Containers:   []NewContainer{{}},
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

func TestContainerRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       ContainerRequest
		wantErr bool
	}{
		{"valid", ContainerRequest{OrderRequest: ValidOrderRequest(), ContainerID: "c1"}, false},
		{"bad_order", ContainerRequest{ContainerID: "c1"}, true},
		{"missing_container_id", ContainerRequest{OrderRequest: ValidOrderRequest()}, true},
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

func TestSetContainerFieldsRequest_Validate(t *testing.T) {
	validReq := ContainerRequest{OrderRequest: ValidOrderRequest(), ContainerID: "c1"}
	tests := []struct {
		name    string
		v       SetContainerFieldsRequest
		wantErr bool
	}{
		{"valid", SetContainerFieldsRequest{
			ContainerRequest: validReq,
			SetFieldsRequest: SetFieldsRequest{SetStrings: map[string]string{"number": "ABC"}},
		}, false},
		{"valid_no_fields", SetContainerFieldsRequest{ContainerRequest: validReq}, false},
		{"bad_container_request", SetContainerFieldsRequest{
			SetFieldsRequest: SetFieldsRequest{SetStrings: map[string]string{"number": "ABC"}},
		}, true},
		{"untrimmed_value", SetContainerFieldsRequest{
			ContainerRequest: validReq,
			SetFieldsRequest: SetFieldsRequest{SetStrings: map[string]string{"number": " ABC "}},
		}, true},
		{"empty_string_name", SetContainerFieldsRequest{
			ContainerRequest: validReq,
			SetFieldsRequest: SetFieldsRequest{SetStrings: map[string]string{"": "ABC"}},
		}, true},
		{"unknown_string_name", SetContainerFieldsRequest{
			ContainerRequest: validReq,
			SetFieldsRequest: SetFieldsRequest{SetStrings: map[string]string{"bogus": "ABC"}},
		}, true},
		{"unknown_date_name", SetContainerFieldsRequest{
			ContainerRequest: validReq,
			SetFieldsRequest: SetFieldsRequest{SetDates: map[string]string{"someDate": "2023-01-01"}},
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
