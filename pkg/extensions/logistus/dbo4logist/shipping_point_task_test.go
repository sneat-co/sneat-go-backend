package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateShippingPointTasksRequest(t *testing.T) {
	assert.NoError(t, ValidateShippingPointTasksRequest([]ShippingPointTask{ShippingPointTaskPick}, true))
	assert.NoError(t, ValidateShippingPointTasksRequest(nil, false))
	assert.Error(t, ValidateShippingPointTasksRequest(nil, true))
	assert.Error(t, ValidateShippingPointTasksRequest([]ShippingPointTask{"bad"}, false))
}

func TestValidateShippingPointTasksRecord(t *testing.T) {
	assert.NoError(t, ValidateShippingPointTasksRecord([]ShippingPointTask{ShippingPointTaskLoad}))
	assert.Error(t, ValidateShippingPointTasksRecord(nil))
	assert.Error(t, ValidateShippingPointTasksRecord([]ShippingPointTask{"bad"}))
}
