package dto4logist

import (
	//"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCreateOrderRequest_Validate validates CreateOrderRequest.Validate() method
func TestCreateOrderRequest_Validate(t *testing.T) {
	type fields struct {
		TeamRequest        dto4teamus.TeamRequest
		Order              models4logist.OrderBase
		NumberOfContainers map[string]int
	}
	tests := []struct {
		name        string
		fields      fields
		expectedErr string
	}{
		{name: "empty", fields: fields{}, expectedErr: "validation error: invalid request: bad value for field [team]: missing required field"},
		{name: "should_pass", fields: fields{
			TeamRequest: dto4teamus.TeamRequest{TeamID: "teamID"},
			Order: models4logist.OrderBase{
				Direction: models4logist.OrderDirectionExport,
				Status:    models4logist.OrderStatusDraft,
				Route: &models4logist.OrderRoute{
					Origin:      models4logist.TransitPoint{CountryID: "IE"},
					Destination: models4logist.TransitPoint{CountryID: "UK"},
				},
			},
			NumberOfContainers: map[string]int{"containerID": 1},
		}, expectedErr: ""},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := CreateOrderRequest{
				TeamRequest:        tt.fields.TeamRequest,
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
