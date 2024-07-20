package dto4logist

import (
	//"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCreateOrderRequest_Validate validates CreateOrderRequest.Validate() method
func TestCreateOrderRequest_Validate(t *testing.T) {
	type fields struct {
		SpaceRequest       dto4spaceus.SpaceRequest
		Order              dbo4logist.OrderBase
		NumberOfContainers map[string]int
	}
	tests := []struct {
		name        string
		fields      fields
		expectedErr string
	}{
		{name: "empty", fields: fields{}, expectedErr: "validation error: invalid request: bad value for field [space]: missing required field"},
		{name: "should_pass", fields: fields{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: "space1"},
			Order: dbo4logist.OrderBase{
				Direction: dbo4logist.OrderDirectionExport,
				Status:    dbo4logist.OrderStatusDraft,
				Route: &dbo4logist.OrderRoute{
					Origin:      dbo4logist.TransitPoint{CountryID: "IE"},
					Destination: dbo4logist.TransitPoint{CountryID: "UK"},
				},
			},
			NumberOfContainers: map[string]int{"containerID": 1},
		}, expectedErr: ""},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := CreateOrderRequest{
				SpaceRequest:       tt.fields.SpaceRequest,
				Order:              tt.fields.Order,
				NumberOfContainers: tt.fields.NumberOfContainers,
			}
			err := v.Validate()
			if err != nil {
				if tt.expectedErr == "" {
					t.Errorf("Unexpected error: %v", err)
				} else if errMessage := err.Error(); errMessage != tt.expectedErr {
					assert.Equal(t, tt.expectedErr, errMessage, "Unexpected error")
				}
			} else if tt.expectedErr != "" {
				t.Errorf("Passed validation but expected error: %v", tt.expectedErr)
			}
		})
	}
}
