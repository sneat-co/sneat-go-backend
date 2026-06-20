package dbo4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func TestValidateOrderShippingPointStatus(t *testing.T) {
	assert.NoError(t, validateOrderShippingPointStatus("status", OrderShippingPointStatusPending))
	assert.NoError(t, validateOrderShippingPointStatus("status", OrderShippingPointStatusProcessing))
	assert.NoError(t, validateOrderShippingPointStatus("status", OrderShippingPointStatusCompleted))
	assert.Error(t, validateOrderShippingPointStatus("status", ""))
	assert.Error(t, validateOrderShippingPointStatus("status", "bad"))
}

func TestShippingPointLocation_Validate(t *testing.T) {
	validAddr := &dbmodels.Address{CountryID: "US"}
	tests := []struct {
		name    string
		v       ShippingPointLocation
		wantErr bool
	}{
		{"valid", ShippingPointLocation{ContactID: "c1", Title: "T1", Address: validAddr}, false},
		{"missing_contact", ShippingPointLocation{Title: "T1", Address: validAddr}, true},
		{"missing_title", ShippingPointLocation{ContactID: "c1", Address: validAddr}, true},
		{"missing_address", ShippingPointLocation{ContactID: "c1", Title: "T1"}, true},
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

func validShippingPoint(id string) *OrderShippingPoint {
	return &OrderShippingPoint{
		ID: id,
		ShippingPointBase: ShippingPointBase{
			Status:       ShippingPointStatusPending,
			FreightPoint: FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskPick}},
		},
		Counterparty: ShippingPointCounterparty{ContactID: "cp-" + id, Title: "T"},
	}
}

func TestOrderShippingPoint_Validate_table(t *testing.T) {
	addr := &dbmodels.Address{CountryID: "US"}
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, validShippingPoint("sp1").Validate())
	})
	t.Run("missing_id", func(t *testing.T) {
		v := validShippingPoint("sp1")
		v.ID = ""
		assert.Error(t, v.Validate())
	})
	t.Run("missing_tasks", func(t *testing.T) {
		v := validShippingPoint("sp1")
		v.Tasks = nil
		assert.Error(t, v.Validate())
	})
	t.Run("bad_status", func(t *testing.T) {
		v := validShippingPoint("sp1")
		v.Status = "bad"
		assert.Error(t, v.Validate())
	})
	t.Run("bad_counterparty", func(t *testing.T) {
		v := validShippingPoint("sp1")
		v.Counterparty = ShippingPointCounterparty{}
		assert.Error(t, v.Validate())
	})
	t.Run("location_same_contact_as_counterparty", func(t *testing.T) {
		v := validShippingPoint("sp1")
		v.Location = &ShippingPointLocation{ContactID: v.Counterparty.ContactID, Title: "L", Address: addr}
		assert.Error(t, v.Validate())
	})
	t.Run("valid_with_location", func(t *testing.T) {
		v := validShippingPoint("sp1")
		v.Location = &ShippingPointLocation{ContactID: "loc1", Title: "L", Address: addr}
		assert.NoError(t, v.Validate())
	})
}

func TestWithShippingPoints_Validate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.NoError(t, (&WithShippingPoints{}).Validate())
	})
	t.Run("valid", func(t *testing.T) {
		v := &WithShippingPoints{ShippingPoints: []*OrderShippingPoint{validShippingPoint("sp1"), validShippingPoint("sp2")}}
		assert.NoError(t, v.Validate())
	})
	t.Run("invalid_child", func(t *testing.T) {
		bad := validShippingPoint("sp1")
		bad.ID = ""
		v := &WithShippingPoints{ShippingPoints: []*OrderShippingPoint{bad}}
		assert.Error(t, v.Validate())
	})
	t.Run("duplicate_id", func(t *testing.T) {
		v := &WithShippingPoints{ShippingPoints: []*OrderShippingPoint{validShippingPoint("sp1"), validShippingPoint("sp1")}}
		assert.Error(t, v.Validate())
	})
}

func TestWithShippingPoints_Updates(t *testing.T) {
	empty := (&WithShippingPoints{}).Updates()
	assert.Len(t, empty, 1)
	assert.Equal(t, "shippingPoints", empty[0].FieldName())

	withData := (&WithShippingPoints{ShippingPoints: []*OrderShippingPoint{validShippingPoint("sp1")}}).Updates()
	assert.Len(t, withData, 1)
	assert.Equal(t, "shippingPoints", withData[0].FieldName())
}

func TestWithShippingPoints_GetShippingPointByID(t *testing.T) {
	sp := validShippingPoint("sp1")
	v := &WithShippingPoints{ShippingPoints: []*OrderShippingPoint{sp}}
	i, got := v.GetShippingPointByID("sp1")
	assert.Equal(t, 0, i)
	assert.Equal(t, sp, got)
	i, got = v.GetShippingPointByID("missing")
	assert.Equal(t, -1, i)
	assert.Nil(t, got)
}

func TestWithShippingPoints_GetShippingPointByContactID(t *testing.T) {
	sp := validShippingPoint("sp1") // counterparty contactID = cp-sp1
	sp.Location = &ShippingPointLocation{ContactID: "loc1", Title: "L", Address: &dbmodels.Address{CountryID: "US"}}
	v := &WithShippingPoints{ShippingPoints: []*OrderShippingPoint{sp}}

	i, got := v.GetShippingPointByContactID("cp-sp1")
	assert.Equal(t, 0, i)
	assert.Equal(t, sp, got)

	i, got = v.GetShippingPointByContactID("loc1")
	assert.Equal(t, 0, i)
	assert.Equal(t, sp, got)

	i, got = v.GetShippingPointByContactID("missing")
	assert.Equal(t, -1, i)
	assert.Nil(t, got)
}
