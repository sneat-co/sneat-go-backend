package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderContainerBase_String(t *testing.T) {
	v := OrderContainerBase{Type: ContainerType20ft, Number: "ABC123"}
	s := v.String()
	assert.Contains(t, s, "20ft")
	assert.Contains(t, s, "ABC123")
}

func TestOrderContainerBase_Validate_table(t *testing.T) {
	tests := []struct {
		name    string
		v       OrderContainerBase
		wantErr bool
	}{
		{"valid", OrderContainerBase{Type: ContainerType20ft}, false},
		{"unknown_type_keyword", OrderContainerBase{Type: "unknown"}, false},
		{"empty_type", OrderContainerBase{Type: ""}, true},
		{"type_with_spaces", OrderContainerBase{Type: " 20ft"}, true},
		{"unsupported_type", OrderContainerBase{Type: "99ft"}, true},
		{"number_with_spaces", OrderContainerBase{Type: ContainerType20ft, Number: " 1 "}, true},
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

func TestOrderContainer_String(t *testing.T) {
	v := OrderContainer{ID: "C1", OrderContainerBase: OrderContainerBase{Type: ContainerType20ft, Number: "N1"}}
	s := v.String()
	assert.Contains(t, s, "C1")
	assert.Contains(t, s, "20ft")
	assert.Contains(t, s, "N1")
}

func TestValidateContainerID(t *testing.T) {
	assert.NoError(t, validateContainerID("id", "c1"))
	assert.Error(t, validateContainerID("id", ""))
	assert.Error(t, validateContainerID("id", "  "))
	assert.Error(t, validateContainerID("id", "bad id with spaces"))
}

func TestOrderContainer_Validate_table(t *testing.T) {
	tests := []struct {
		name    string
		v       OrderContainer
		wantErr bool
	}{
		{"valid", OrderContainer{ID: "c1", OrderContainerBase: OrderContainerBase{Type: ContainerType20ft}}, false},
		{"missing_id", OrderContainer{OrderContainerBase: OrderContainerBase{Type: ContainerType20ft}}, true},
		{"bad_base", OrderContainer{ID: "c1", OrderContainerBase: OrderContainerBase{Type: "bad"}}, true},
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

func TestWithOrderContainers_Validate_table(t *testing.T) {
	tests := []struct {
		name    string
		v       WithOrderContainers
		wantErr bool
	}{
		{"empty", WithOrderContainers{}, false},
		{"valid", WithOrderContainers{Containers: []*OrderContainer{
			{ID: "c1", OrderContainerBase: OrderContainerBase{Type: ContainerType20ft, Number: "N1"}},
			{ID: "c2", OrderContainerBase: OrderContainerBase{Type: ContainerType40ft, Number: "N2"}},
		}}, false},
		{"invalid_child", WithOrderContainers{Containers: []*OrderContainer{
			{ID: "", OrderContainerBase: OrderContainerBase{Type: ContainerType20ft}},
		}}, true},
		{"duplicate_id", WithOrderContainers{Containers: []*OrderContainer{
			{ID: "c1", OrderContainerBase: OrderContainerBase{Type: ContainerType20ft}},
			{ID: "c1", OrderContainerBase: OrderContainerBase{Type: ContainerType40ft}},
		}}, true},
		{"duplicate_number", WithOrderContainers{Containers: []*OrderContainer{
			{ID: "c1", OrderContainerBase: OrderContainerBase{Type: ContainerType20ft, Number: "N1"}},
			{ID: "c2", OrderContainerBase: OrderContainerBase{Type: ContainerType40ft, Number: "N1"}},
		}}, true},
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

func TestWithOrderContainers_GenerateRandomContainerID(t *testing.T) {
	v := WithOrderContainers{Containers: []*OrderContainer{{ID: "c1"}}}
	id := v.GenerateRandomContainerID()
	assert.NotEmpty(t, id)
}

func TestWithOrderContainers_Updates(t *testing.T) {
	empty := WithOrderContainers{}.Updates()
	assert.Len(t, empty, 1)
	assert.Equal(t, "containers", empty[0].FieldName())

	withData := WithOrderContainers{Containers: []*OrderContainer{{ID: "c1"}}}.Updates()
	assert.Len(t, withData, 1)
	assert.Equal(t, "containers", withData[0].FieldName())
}
