package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerEndpoint_IsEmpty(t *testing.T) {
	assert.True(t, ContainerEndpoint{}.IsEmpty())
	assert.False(t, ContainerEndpoint{ScheduledDate: "2023-01-01"}.IsEmpty())
	assert.False(t, ContainerEndpoint{ActualDate: "2023-01-01"}.IsEmpty())
	assert.False(t, ContainerEndpoint{ByContactID: "c1"}.IsEmpty())
}

func TestContainerEndpoint_String(t *testing.T) {
	s := ContainerEndpoint{ByContactID: "c1", ScheduledDate: "2023-01-01", ActualDate: "2023-01-02"}.String()
	assert.Contains(t, s, "c1")
	assert.Contains(t, s, "2023-01-01")
	assert.Contains(t, s, "2023-01-02")
}

func TestContainerEndpoint_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       ContainerEndpoint
		wantErr bool
	}{
		{"empty", ContainerEndpoint{}, false},
		{"valid", ContainerEndpoint{ScheduledDate: "2023-01-01", ActualDate: "2023-01-02", ByContactID: "c1"}, false},
		{"bad_scheduled", ContainerEndpoint{ScheduledDate: "bad"}, true},
		{"bad_actual", ContainerEndpoint{ActualDate: "bad"}, true},
		{"contact_with_spaces", ContainerEndpoint{ByContactID: " c1 "}, true},
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

func TestContainerEndpoints_Validate(t *testing.T) {
	assert.NoError(t, ContainerEndpoints{}.Validate())
	assert.NoError(t, ContainerEndpoints{
		Arrival:   &ContainerEndpoint{ScheduledDate: "2023-01-01"},
		Departure: &ContainerEndpoint{ScheduledDate: "2023-01-02"},
	}.Validate())
	assert.Error(t, ContainerEndpoints{Arrival: &ContainerEndpoint{ScheduledDate: "bad"}}.Validate())
	assert.Error(t, ContainerEndpoints{Departure: &ContainerEndpoint{ScheduledDate: "bad"}}.Validate())
}

func TestContainerEndpoints_Strings(t *testing.T) {
	s := ContainerEndpoints{
		Arrival:   &ContainerEndpoint{ScheduledDate: "2023-01-01"},
		Departure: &ContainerEndpoint{ScheduledDate: "2023-01-02"},
	}.Strings()
	assert.Contains(t, s, "Arrival")
	assert.Contains(t, s, "Departure")
}
