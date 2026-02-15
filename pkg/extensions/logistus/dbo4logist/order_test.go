package dbo4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func TestOrderCounterparty_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       OrderCounterparty
		wantErr bool
	}{
		{"valid", OrderCounterparty{ContactID: "c1", Role: CounterpartyRoleBuyer, CountryID: "US", Title: "Buyer 1"}, false},
		{"missing_contact", OrderCounterparty{Role: CounterpartyRoleBuyer, CountryID: "US", Title: "Buyer 1"}, true},
		{"missing_country", OrderCounterparty{ContactID: "c1", Role: CounterpartyRoleBuyer, Title: "Buyer 1"}, true},
		{"missing_title", OrderCounterparty{ContactID: "c1", Role: CounterpartyRoleBuyer, CountryID: "US"}, true},
		{"invalid_role", OrderCounterparty{ContactID: "c1", Role: "invalid", CountryID: "US", Title: "Buyer 1"}, true},
		{"same_parent", OrderCounterparty{ContactID: "c1", Role: CounterpartyRoleBuyer, CountryID: "US", Title: "Buyer 1", Parent: &CounterpartyParent{ContactID: "c1"}}, true},
		{"ship_no_country", OrderCounterparty{ContactID: "c1", Role: CounterpartyRoleShip, Title: "Ship 1"}, false},
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

func TestWithCounterparties_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       WithCounterparties
		wantErr bool
	}{
		{"empty", WithCounterparties{}, false},
		{"valid", WithCounterparties{Counterparties: []*OrderCounterparty{
			{ContactID: "c1", Role: CounterpartyRoleBuyer, CountryID: "US", Title: "Buyer 1"},
			{ContactID: "c2", Role: CounterpartyRoleReceiver, CountryID: "US", Title: "Receiver 1"},
		}}, false},
		{"duplicate_buyer", WithCounterparties{Counterparties: []*OrderCounterparty{
			{ContactID: "c1", Role: CounterpartyRoleBuyer, CountryID: "US", Title: "Buyer 1"},
			{ContactID: "c2", Role: CounterpartyRoleBuyer, CountryID: "US", Title: "Buyer 2"},
		}}, true},
		{"allowed_duplicate_trucker", WithCounterparties{Counterparties: []*OrderCounterparty{
			{ContactID: "c1", Role: CounterpartyRoleTrucker, CountryID: "US", Title: "Trucker 1"},
			{ContactID: "c2", Role: CounterpartyRoleTrucker, CountryID: "US", Title: "Trucker 2"},
		}}, false},
		{"same_port", WithCounterparties{Counterparties: []*OrderCounterparty{
			{ContactID: "p1", Role: CounterpartyRolePortFrom, CountryID: "US", Title: "Port"},
			{ContactID: "p1", Role: CounterpartyRolePortTo, CountryID: "US", Title: "Port"},
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

func TestOrderRoute_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       OrderRoute
		wantErr bool
	}{
		{"valid", OrderRoute{
			Origin:      TransitPoint{CountryID: "US"},
			Destination: TransitPoint{CountryID: "GB"},
		}, false},
		{"invalid_origin", OrderRoute{
			Destination: TransitPoint{CountryID: "GB"},
		}, true},
		{"invalid_destination", OrderRoute{
			Origin: TransitPoint{CountryID: "US"},
		}, true},
		{"invalid_transit", OrderRoute{
			Origin:        TransitPoint{CountryID: "US"},
			Destination:   TransitPoint{CountryID: "GB"},
			TransitPoints: []*TransitPoint{{CountryID: ""}},
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

func TestOrderStep_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       OrderStep
		wantErr bool
	}{
		{"valid_pending", OrderStep{ID: "s1", Status: OrderStepStatusPending}, false},
		{"valid_completed", OrderStep{ID: "s1", Status: OrderStepStatusCompleted}, false},
		{"missing_id", OrderStep{Status: OrderStepStatusPending}, true},
		{"missing_status", OrderStep{ID: "s1"}, true},
		{"invalid_status", OrderStep{ID: "s1", Status: "invalid"}, true},
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

func TestOrderDbo_Validate(t *testing.T) {
	order := &OrderDbo{
		WithSpaceID:  dbmodels.WithSpaceID{SpaceID: coretypes.SpaceID("s1")},
		WithSpaceIDs: dbmodels.WithSpaceIDs{SpaceIDs: []coretypes.SpaceID{coretypes.SpaceID("s1")}},
		OrderBase: OrderBase{
			Status:    OrderStatusDraft,
			Direction: OrderDirectionImport,
			Route: &OrderRoute{
				Origin:      TransitPoint{CountryID: "US"},
				Destination: TransitPoint{CountryID: "GB"},
			},
		},
	}
	// Many sub-validators will fail if not initialized
	err := order.Validate()
	assert.Error(t, err) // Expecting some errors due to many required fields in sub-structs
}

func TestFreightLoad_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       FreightLoad
		wantErr bool
	}{
		{"valid", FreightLoad{NumberOfPallets: 1, GrossWeightKg: 100, VolumeM3: 10}, false},
		{"negative_pallets", FreightLoad{NumberOfPallets: -1}, true},
		{"negative_weight", FreightLoad{GrossWeightKg: -1}, true},
		{"negative_volume", FreightLoad{VolumeM3: -1}, true},
		{"untrimmed_note", FreightLoad{Note: " note "}, true},
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

func TestContainerPoint_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       ContainerPoint
		wantErr bool
	}{
		{"valid", ContainerPoint{
			ContainerID:     "c1",
			ShippingPointID: "sp1",
			ShippingPointBase: ShippingPointBase{
				Status: ShippingPointStatusPending,
				FreightPoint: FreightPoint{
					Tasks: []ShippingPointTask{ShippingPointTaskPick},
				},
			},
		}, false},
		{"missing_sp_id", ContainerPoint{ContainerID: "c1"}, true},
		{"missing_c_id", ContainerPoint{ShippingPointID: "sp1"}, true},
		{"invalid_status", ContainerPoint{ContainerID: "c1", ShippingPointID: "sp1", ShippingPointBase: ShippingPointBase{Status: "invalid"}}, true},
		{"ref_too_long", ContainerPoint{
			ContainerID:     "c1",
			ShippingPointID: "sp1",
			RefNumber:       "this is a very long reference number that exceeds fifty characters limit",
			ShippingPointBase: ShippingPointBase{
				Status: ShippingPointStatusPending,
				FreightPoint: FreightPoint{
					Tasks: []ShippingPointTask{ShippingPointTaskPick},
				},
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

func TestShippingPointBase_String(t *testing.T) {
	v := ShippingPointBase{Status: ShippingPointStatusPending}
	s := v.String()
	assert.Contains(t, s, "status=pending")
}

func TestContainerPoint_String(t *testing.T) {
	v := ContainerPoint{ContainerID: "C1", ShippingPointID: "SP1", ShippingPointBase: ShippingPointBase{Status: ShippingPointStatusPending}}
	s := v.String()
	assert.Contains(t, s, "C1")
	assert.Contains(t, s, "SP1")
}

func TestFreightLoad_IsEmpty(t *testing.T) {
	var v *FreightLoad
	assert.True(t, v.IsEmpty())
	v = &FreightLoad{}
	assert.True(t, v.IsEmpty())
	v = &FreightLoad{NumberOfPallets: 1}
	assert.False(t, v.IsEmpty())
}

func TestFreightLoad_Add(t *testing.T) {
	v1 := &FreightLoad{NumberOfPallets: 1, GrossWeightKg: 10, VolumeM3: 5}
	v2 := &FreightLoad{NumberOfPallets: 2, GrossWeightKg: 20, VolumeM3: 10}
	v1.Add(v2)
	assert.Equal(t, 3, v1.NumberOfPallets)
	assert.Equal(t, 30, v1.GrossWeightKg)
	assert.Equal(t, 15, v1.VolumeM3)
	v1.Add(nil)
	assert.Equal(t, 3, v1.NumberOfPallets)
}

func TestShippingPointTask_Validate(t *testing.T) {
	assert.NoError(t, ValidateShippingPointTask(ShippingPointTaskPick, ValidatingRecord, func() string { return "f" }))
	assert.Error(t, ValidateShippingPointTask("invalid", ValidatingRecord, func() string { return "f" }))
	assert.Error(t, ValidateShippingPointTask(" pick", ValidatingRecord, func() string { return "f" }))
	assert.Error(t, ValidateShippingPointTask("", ValidatingRecord, func() string { return "f" }))
	assert.Error(t, ValidateShippingPointTask("", ValidatingRequest, func() string { return "f" }))
}

func TestShippingPointOrder_validateOrderShippingPointStatus(t *testing.T) {
	assert.NoError(t, validateOrderShippingPointStatus("f", "pending"))
	assert.Error(t, validateOrderShippingPointStatus("f", "invalid"))
	assert.Error(t, validateOrderShippingPointStatus("f", ""))
}

func TestShippingPointCounterparty_Validate(t *testing.T) {
	v := ShippingPointCounterparty{ContactID: "c1", Title: "T1"}
	assert.NoError(t, v.Validate())
	v = ShippingPointCounterparty{Title: "T1"}
	assert.Error(t, v.Validate())
	v = ShippingPointCounterparty{ContactID: "c1"}
	assert.Error(t, v.Validate())
}

func TestOrderShippingPoint_Validate(t *testing.T) {
	v := OrderShippingPoint{
		ID: "sp1",
		ShippingPointBase: ShippingPointBase{
			Status: ShippingPointStatusPending,
			FreightPoint: FreightPoint{
				Tasks: []ShippingPointTask{ShippingPointTaskPick},
			},
		},
		Counterparty: ShippingPointCounterparty{ContactID: "c1", Title: "T1"},
	}
	assert.NoError(t, v.Validate())
	v.ID = ""
	assert.Error(t, v.Validate())
}

func TestGetCounterpartyTypeByRole(t *testing.T) {
	assert.Equal(t, CounterpartyTypeLocation, GetCounterpartyTypeByRole(CounterpartyRolePickPoint))
	assert.Equal(t, CounterpartyTypeCompany, GetCounterpartyTypeByRole(CounterpartyRoleDispatcher))
	assert.Equal(t, CounterpartyTypeUnknown, GetCounterpartyTypeByRole("unknown"))
}

func TestOrderCounterparty_String(t *testing.T) {
	v := OrderCounterparty{ContactID: "c1", Role: "r1", CountryID: "US", Title: "t1"}
	assert.Contains(t, v.String(), "c1")
}

func TestOrderContainerBase_Validate(t *testing.T) {
	v := OrderContainerBase{Type: ContainerType20ft}
	assert.NoError(t, v.Validate())
	v.Type = ""
	assert.Error(t, v.Validate())
	v.Type = ContainerType20ft
	v.Instructions = "some instructions"
	assert.NoError(t, v.Validate())
}

func TestOrderContainer_Validate(t *testing.T) {
	v := OrderContainer{ID: "C1", OrderContainerBase: OrderContainerBase{Type: ContainerType20ft}}
	assert.NoError(t, v.Validate())
	v.ID = ""
	assert.Error(t, v.Validate())
}

func TestWithOrderContainers_GetContainerIDs(t *testing.T) {
	v := WithOrderContainers{Containers: []*OrderContainer{{ID: "C1"}, {ID: "C2"}}}
	ids := v.GetContainerIDs()
	assert.ElementsMatch(t, []string{"C1", "C2"}, ids)
}

func TestWithOrderContainers_GetContainerByID(t *testing.T) {
	v := WithOrderContainers{Containers: []*OrderContainer{{ID: "C1"}}}
	i, c := v.GetContainerByID("C1")
	assert.Equal(t, 0, i)
	assert.NotNil(t, c)
	i, c = v.GetContainerByID("C2")
	assert.Equal(t, -1, i)
	assert.Nil(t, c)
}

func TestWithOrderContainers_RemoveContainer(t *testing.T) {
	v := WithOrderContainers{Containers: []*OrderContainer{{ID: "C1"}, {ID: "C2"}}}
	newContainers, found := v.RemoveContainer("C1")
	assert.True(t, found)
	assert.Len(t, newContainers, 1)
	assert.Equal(t, "C2", newContainers[0].ID)
}

func TestOrderBase_Validate(t *testing.T) {
	v := OrderBase{
		Status:    OrderStatusDraft,
		Direction: OrderDirectionImport,
		Route: &OrderRoute{
			Origin:      TransitPoint{CountryID: "US"},
			Destination: TransitPoint{CountryID: "GB"},
		},
	}
	assert.NoError(t, v.Validate())
	v.Route = nil
	assert.NoError(t, v.Validate()) // Route is optional for draft
}

func TestSegment_Validate(t *testing.T) {
	v := SegmentDates{Arrives: "2023-01-01"}
	assert.NoError(t, v.Validate())
}

func TestSegmentCounterparty_Validate(t *testing.T) {
	v := SegmentCounterparty{ContactID: "c1", Role: CounterpartyRoleTrucker}
	assert.NoError(t, v.Validate())
	v.Role = ""
	assert.Error(t, v.Validate())
}

func TestSegmentEndpoint_Validate(t *testing.T) {
	v := SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "c1", Role: CounterpartyRoleTrucker}}
	assert.NoError(t, v.Validate())
	v.Role = CounterpartyRoleDispatchPoint
	assert.Error(t, v.Validate()) // missing ShippingPointID
}

func TestContainerSegment_Validate(t *testing.T) {
	v := &ContainerSegment{}
	v.ContainerID = "c1"
	v.From = SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "c1", Role: CounterpartyRoleTrucker}}
	v.To = SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "c2", Role: CounterpartyRoleTrucker}}
	assert.NoError(t, v.Validate())
}
