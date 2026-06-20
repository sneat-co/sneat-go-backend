package dto4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/stretchr/testify/assert"
)

// validContainerPointRequest returns a ContainerPointRequest that passes Validate().
func validContainerPointRequest() ContainerPointRequest {
	return ContainerPointRequest{
		OrderRequest:    ValidOrderRequest(),
		ShippingPointID: "sp1",
		ContainerID:     "c1",
	}
}

func TestContainerPointsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       ContainerPointsRequest
		wantErr bool
	}{
		{"valid", ContainerPointsRequest{
			OrderRequest:     ValidOrderRequest(),
			ContainerID:      "c1",
			ShippingPointIDs: []string{"sp1"},
		}, false},
		{"bad_order", ContainerPointsRequest{
			ContainerID:      "c1",
			ShippingPointIDs: []string{"sp1"},
		}, true},
		{"missing_container_id", ContainerPointsRequest{
			OrderRequest:     ValidOrderRequest(),
			ShippingPointIDs: []string{"sp1"},
		}, true},
		{"empty_shipping_points", ContainerPointsRequest{
			OrderRequest: ValidOrderRequest(),
			ContainerID:  "c1",
		}, true},
		{"bad_shipping_point", ContainerPointsRequest{
			OrderRequest:     ValidOrderRequest(),
			ContainerID:      "c1",
			ShippingPointIDs: []string{" bad "},
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

func TestContainerPointRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       ContainerPointRequest
		wantErr bool
	}{
		{"valid", validContainerPointRequest(), false},
		{"bad_order", ContainerPointRequest{ShippingPointID: "sp1", ContainerID: "c1"}, true},
		{"missing_container_id", ContainerPointRequest{OrderRequest: ValidOrderRequest(), ShippingPointID: "sp1"}, true},
		{"bad_shipping_point", ContainerPointRequest{OrderRequest: ValidOrderRequest(), ContainerID: "c1", ShippingPointID: " bad "}, true},
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

func validContainerPoint() dbo4logist.ContainerPoint {
	return dbo4logist.ContainerPoint{
		ContainerID:     "c1",
		ShippingPointID: "sp1",
		ShippingPointBase: dbo4logist.ShippingPointBase{
			Status: dbo4logist.ShippingPointStatusPending,
			FreightPoint: dbo4logist.FreightPoint{
				Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskPick},
			},
		},
	}
}

func TestAddContainerPointsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       AddContainerPointsRequest
		wantErr bool
	}{
		{"valid", AddContainerPointsRequest{
			OrderRequest:    ValidOrderRequest(),
			ContainerPoints: []dbo4logist.ContainerPoint{validContainerPoint()},
		}, false},
		{"valid_empty_points", AddContainerPointsRequest{
			OrderRequest: ValidOrderRequest(),
		}, false},
		{"bad_order", AddContainerPointsRequest{
			ContainerPoints: []dbo4logist.ContainerPoint{validContainerPoint()},
		}, true},
		{"bad_container_point", AddContainerPointsRequest{
			OrderRequest:    ValidOrderRequest(),
			ContainerPoints: []dbo4logist.ContainerPoint{{ContainerID: "c1"}},
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

func TestUpdateContainerPointRequest_Validate(t *testing.T) {
	valid := UpdateContainerPointRequest{
		ContainerPointRequest: validContainerPointRequest(),
		FreightPoint: dbo4logist.FreightPoint{
			Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskPick},
		},
	}
	tests := []struct {
		name    string
		v       UpdateContainerPointRequest
		wantErr bool
	}{
		{"valid", valid, false},
		{"bad_order", UpdateContainerPointRequest{
			ContainerPointRequest: ContainerPointRequest{ShippingPointID: "sp1", ContainerID: "c1"},
		}, true},
		{"bad_shipping_point", UpdateContainerPointRequest{
			ContainerPointRequest: ContainerPointRequest{OrderRequest: ValidOrderRequest(), ContainerID: "c1", ShippingPointID: " bad "},
		}, true},
		{"missing_container_id", UpdateContainerPointRequest{
			ContainerPointRequest: ContainerPointRequest{OrderRequest: ValidOrderRequest(), ShippingPointID: "sp1"},
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

func TestSetContainerPointTaskRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       SetContainerPointTaskRequest
		wantErr bool
	}{
		{"valid", SetContainerPointTaskRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Task:                  dbo4logist.ShippingPointTaskPick,
		}, false},
		{"bad_request", SetContainerPointTaskRequest{
			Task: dbo4logist.ShippingPointTaskPick,
		}, true},
		{"bad_task", SetContainerPointTaskRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Task:                  "invalid",
		}, true},
		{"empty_task", SetContainerPointTaskRequest{
			ContainerPointRequest: validContainerPointRequest(),
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

func TestSetContainerEndpointFieldsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       SetContainerEndpointFieldsRequest
		wantErr bool
	}{
		{"valid_dates", SetContainerEndpointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Side:                  dbo4logist.EndpointSideArrival,
			Dates:                 map[string]string{"eta": "2023-01-01"},
		}, false},
		{"valid_times", SetContainerEndpointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Side:                  dbo4logist.EndpointSideDeparture,
			Times:                 map[string]string{"etd": "10:00"},
		}, false},
		{"bad_request", SetContainerEndpointFieldsRequest{
			Side:  dbo4logist.EndpointSideArrival,
			Dates: map[string]string{"eta": "2023-01-01"},
		}, true},
		{"missing_side", SetContainerEndpointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Dates:                 map[string]string{"eta": "2023-01-01"},
		}, true},
		{"bad_side", SetContainerEndpointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Side:                  "sideways",
			Dates:                 map[string]string{"eta": "2023-01-01"},
		}, true},
		{"nothing_set", SetContainerEndpointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Side:                  dbo4logist.EndpointSideArrival,
		}, true},
		{"empty_date_name", SetContainerEndpointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Side:                  dbo4logist.EndpointSideArrival,
			Dates:                 map[string]string{"": "2023-01-01"},
		}, true},
		{"bad_date_value", SetContainerEndpointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Side:                  dbo4logist.EndpointSideArrival,
			Dates:                 map[string]string{"eta": "not-a-date"},
		}, true},
		{"empty_time_name", SetContainerEndpointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Side:                  dbo4logist.EndpointSideArrival,
			Times:                 map[string]string{"": "10:00"},
		}, true},
		{"bad_time_value", SetContainerEndpointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Side:                  dbo4logist.EndpointSideArrival,
			Times:                 map[string]string{"etd": "not-a-time"},
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

func TestSetContainerPointFieldsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       SetContainerPointFieldsRequest
		wantErr bool
	}{
		{"valid_notes", SetContainerPointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			SetStrings:            map[string]string{"notes": "hello"},
		}, false},
		{"valid_ref", SetContainerPointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			SetStrings:            map[string]string{"refNumber": "ABC123"},
		}, false},
		{"bad_request", SetContainerPointFieldsRequest{
			SetStrings: map[string]string{"notes": "hello"},
		}, true},
		{"empty_name", SetContainerPointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			SetStrings:            map[string]string{"": "hello"},
		}, true},
		{"unknown_name", SetContainerPointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			SetStrings:            map[string]string{"bogus": "hello"},
		}, true},
		{"notes_too_long", SetContainerPointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			SetStrings:            map[string]string{"notes": longString(10001)},
		}, true},
		{"ref_too_long", SetContainerPointFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			SetStrings:            map[string]string{"refNumber": longString(51)},
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

func TestSetContainerPointFreightFieldsRequest_Validate(t *testing.T) {
	negative := -1
	positive := 5
	tests := []struct {
		name    string
		v       SetContainerPointFreightFieldsRequest
		wantErr bool
	}{
		{"valid", SetContainerPointFreightFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Task:                  dbo4logist.ShippingPointTaskPick,
			Integers:              map[string]*int{"pallets": &positive},
		}, false},
		{"bad_request", SetContainerPointFreightFieldsRequest{
			Task: dbo4logist.ShippingPointTaskPick,
		}, true},
		{"negative_integer", SetContainerPointFreightFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Task:                  dbo4logist.ShippingPointTaskPick,
			Integers:              map[string]*int{"pallets": &negative},
		}, true},
		{"bad_task", SetContainerPointFreightFieldsRequest{
			ContainerPointRequest: validContainerPointRequest(),
			Task:                  "invalid",
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

func longString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}
