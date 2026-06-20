package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidOrderGraph_isValid(t *testing.T) {
	order := newValidOrderGraph()
	require.NoError(t, order.Validate())
}

// TestOrderDbo_Validate_mutations mutates a valid graph in ways that should
// cause Validate() to fail, exercising the cross-entity validators.
func TestOrderDbo_Validate_mutations(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(o *OrderDbo)
	}{
		{"dangling_container_point_container", func(o *OrderDbo) {
			o.ContainerPoints[0].ContainerID = "no-such-container"
		}},
		{"dangling_container_point_shipping_point", func(o *OrderDbo) {
			o.ContainerPoints[0].ShippingPointID = "no-such-sp"
		}},
		{"segment_unknown_container", func(o *OrderDbo) {
			// container.validateOrder still passes (container ids unchanged),
			// but WithSegments.validateOrder fails to find the container.
			o.Segments[0].ContainerID = "no-such-container"
		}},
		{"segment_endpoint_counterparty_not_in_order", func(o *OrderDbo) {
			o.Segments[0].From.ContactID = "ghost"
		}},
		{"segment_endpoint_missing_shipping_point", func(o *OrderDbo) {
			o.Segments[0].From.ShippingPointID = "no-such-sp"
		}},
		{"segment_date_mismatch_departure", func(o *OrderDbo) {
			o.Segments[0].Dates.Departs = "2024-12-31"
		}},
		{"segment_date_mismatch_arrival", func(o *OrderDbo) {
			o.Segments[0].Dates.Arrives = "2024-12-31"
		}},
		{"shipping_point_location_contact_not_in_contacts", func(o *OrderDbo) {
			o.ShippingPoints[0].Location.ContactID = "ghost"
		}},
		{"shipping_point_counterparty_role_missing", func(o *OrderDbo) {
			// Remove the dispatcher counterparty referenced by load shipping point.
			cps := o.Counterparties[:0]
			for _, cp := range o.Counterparties {
				if cp.ContactID == "disp" && cp.Role == CounterpartyRoleDispatcher {
					continue
				}
				cps = append(cps, cp)
			}
			o.Counterparties = cps
			o.UpdateKeys()
		}},
		{"dto_shipping_point_no_role_counterparty", func(o *OrderDbo) {
			// Change the load task's location contact to one without a matching
			// dispatch_point counterparty entry.
			o.ShippingPoints[0].Tasks = []ShippingPointTask{ShippingPointTaskLoad}
			o.ShippingPoints[0].Location.ContactID = "recv" // recv has no dispatch_point role
		}},
		{"counterparty_with_unknown_parent", func(o *OrderDbo) {
			o.Counterparties[1].Parent = &CounterpartyParent{ContactID: "ghost", Role: CounterpartyRoleDispatcher}
		}},
		{"counterparty_title_mismatch_contact", func(o *OrderDbo) {
			o.Counterparties[0].Title = "Different Title"
		}},
		{"orphan_contact_not_in_counterparties", func(o *OrderDbo) {
			o.Contacts = append(o.Contacts, &OrderContact{
				ID: "orphan", Type: "company", Title: "Orphan", CountryID: "US",
			})
		}},
		{"missing_origin_country_in_country_ids", func(o *OrderDbo) {
			o.Route = &OrderRoute{
				Origin:      TransitPoint{CountryID: "FR"},
				Destination: TransitPoint{CountryID: "GB"},
			}
		}},
		{"missing_transit_country_in_country_ids", func(o *OrderDbo) {
			o.Route = &OrderRoute{
				Origin:        TransitPoint{CountryID: "US"},
				Destination:   TransitPoint{CountryID: "GB"},
				TransitPoints: []*TransitPoint{{CountryID: "DE"}},
			}
		}},
		{"space_id_not_in_space_ids", func(o *OrderDbo) {
			o.SpaceIDs = nil
		}},
		{"invalid_step", func(o *OrderDbo) {
			o.Steps = []*OrderStep{{ID: "", Status: OrderStepStatusPending}}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := newValidOrderGraph()
			tt.mutate(o)
			assert.Error(t, o.Validate())
		})
	}
}

func TestOrderDbo_UpdateDates(t *testing.T) {
	o := newValidOrderGraph()
	o.Dates = nil
	o.UpdateDates()
	// spd has ScheduledStartDate 2023-01-01, spr has ScheduledEndDate 2023-01-10.
	assert.Contains(t, o.Dates, "2023-01-01")
	assert.Contains(t, o.Dates, "2023-01-10")
	assert.Len(t, o.Dates, 2)

	// Calling again must not duplicate.
	o.UpdateDates()
	assert.Len(t, o.Dates, 2)
}

func TestOrderDbo_updateDatesFromShippingPoint_emptyDates(t *testing.T) {
	o := newValidOrderGraph()
	o.Dates = nil
	o.updateDatesFromShippingPoint(&OrderShippingPoint{}) // no dates -> no-op
	assert.Empty(t, o.Dates)
}

func TestOrderDbo_UpdateCalculatedFields(t *testing.T) {
	o := newValidOrderGraph()
	o.Dates = nil
	o.Keys = nil
	o.UpdateCalculatedFields()
	assert.NotEmpty(t, o.Keys)
	assert.NotEmpty(t, o.Dates)
	// Calculated fields must keep the graph valid.
	require.NoError(t, o.Validate())
}

func TestWithShippingPoints_NewOrderShippingPointID(t *testing.T) {
	v := &WithShippingPoints{ShippingPoints: []*OrderShippingPoint{{ID: "existing"}}}
	id := v.NewOrderShippingPointID()
	assert.NotEmpty(t, id)
	assert.NotEqual(t, "existing", id)
}

func TestWithShippingPoints_DeleteShippingPoint(t *testing.T) {
	v := &WithShippingPoints{ShippingPoints: []*OrderShippingPoint{
		{ID: "spd", Location: &ShippingPointLocation{ContactID: "dp"}},
		{ID: "spr", Location: &ShippingPointLocation{ContactID: "rp"}},
	}}
	deletedID, remaining := v.DeleteShippingPoint("", "dp")
	assert.Equal(t, "spd", deletedID)
	assert.Len(t, remaining, 1)
	assert.Equal(t, "spr", remaining[0].ID)
	assert.Len(t, v.ShippingPoints, 1)

	// Deleting a contact that is not present removes nothing.
	deletedID, remaining = v.DeleteShippingPoint("", "ghost")
	assert.Empty(t, deletedID)
	assert.Len(t, remaining, 1)
}

func TestOrderBrief_Validate(t *testing.T) {
	base := OrderBase{
		Status:    OrderStatusDraft,
		Direction: OrderDirectionImport,
	}
	assert.NoError(t, OrderBrief{ID: "o1", OrderBase: base}.Validate())
	assert.Error(t, OrderBrief{ID: "", OrderBase: base}.Validate())
	assert.Error(t, OrderBrief{ID: "o1", OrderBase: OrderBase{}}.Validate()) // bad OrderBase
}

func TestOrderBase_validateCounterparty_withParent(t *testing.T) {
	// A counterparty with a valid in-order parent passes;
	// an unknown parent fails.
	base := OrderBase{
		Status:    OrderStatusActive,
		Direction: OrderDirectionImport,
		WithCounterparties: WithCounterparties{
			Counterparties: []*OrderCounterparty{
				{ContactID: "disp", Role: CounterpartyRoleDispatcher, CountryID: "US", Title: "Disp"},
				{ContactID: "dp", Role: CounterpartyRoleDispatchPoint, CountryID: "US", Title: "DP",
					Parent: &CounterpartyParent{ContactID: "disp", Role: CounterpartyRoleDispatcher}},
			},
		},
	}
	assert.NoError(t, base.Validate())

	base.Counterparties[1].Parent = &CounterpartyParent{ContactID: "ghost", Role: CounterpartyRoleDispatcher}
	assert.Error(t, base.Validate())
}
