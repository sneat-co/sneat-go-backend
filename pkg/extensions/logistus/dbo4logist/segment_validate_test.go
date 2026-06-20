package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSegmentCounterparty(t *testing.T) {
	order := newValidOrderGraph()
	// Present in order.
	assert.NoError(t, validateSegmentCounterparty(*order,
		SegmentCounterparty{ContactID: "dp", Role: CounterpartyRoleDispatchPoint}))
	// Not present in order.
	assert.Error(t, validateSegmentCounterparty(*order,
		SegmentCounterparty{ContactID: "dp", Role: CounterpartyRoleTrucker}))
}

func TestValidateSegmentEndpoint_portRole(t *testing.T) {
	order := newValidOrderGraph()
	// Add port counterparties so the endpoint counterparty is found in order.
	order.Counterparties = append(order.Counterparties,
		&OrderCounterparty{ContactID: "portF", Role: CounterpartyRolePortFrom, CountryID: "US", Title: "Port From"},
		&OrderCounterparty{ContactID: "portT", Role: CounterpartyRolePortTo, CountryID: "GB", Title: "Port To"},
	)
	seg := &ContainerSegment{ContainerSegmentKey: ContainerSegmentKey{ContainerID: "cont1"}}

	// Port role with no shipping point -> OK.
	ep := SegmentEndpoint{SegmentCounterparty: SegmentCounterparty{ContactID: "portF", Role: CounterpartyRolePortFrom}}
	assert.NoError(t, validateSegmentEndpoint(*order, seg, "from", ep))

	// Port role that wrongly references a shipping point -> error.
	ep.ShippingPointID = "spd"
	assert.Error(t, validateSegmentEndpoint(*order, seg, "from", ep))
}

func TestValidateSegmentEndpoint_unknownField_panics(t *testing.T) {
	order := newValidOrderGraph()
	seg := order.Segments[0]
	ep := seg.From
	assert.Panics(t, func() {
		_ = validateSegmentEndpoint(*order, seg, "middle", ep)
	})
}

func TestValidateOrderSegment_badByContactID(t *testing.T) {
	order := newValidOrderGraph()
	seg := order.Segments[0]
	seg.ByContactID = " spaced " // invalid contact id record field
	assert.Error(t, validateOrderSegment(*order, seg))
}

// TestContainer_validateOrder_freightAccumulation drives the ToLoad/ToUnload
// accumulation branches in WithOrderContainers.validateOrder via Validate().
func TestContainer_validateOrder_freightAccumulation(t *testing.T) {
	order := newValidOrderGraph()
	order.ContainerPoints[0].ToLoad = &FreightLoad{NumberOfPallets: 5, GrossWeightKg: 100, VolumeM3: 10}
	order.ContainerPoints[1].ToUnload = &FreightLoad{NumberOfPallets: 5, GrossWeightKg: 100, VolumeM3: 10}
	// FreightPoint.Validate requires the matching task to be present (already load/unload).
	assert.NoError(t, order.Validate())
}

// TestContainerPoints_validateOrder_byContactID drives the trucker byContactID
// branches in WithContainerPoints.validateOrder.
func TestContainerPoints_validateOrder_byContactID(t *testing.T) {
	t.Run("valid_trucker", func(t *testing.T) {
		order := newValidOrderGraph()
		order.ContainerPoints[0].Departure.ByContactID = "truck"
		order.ContainerPoints[1].Arrival.ByContactID = "truck"
		assert.NoError(t, order.Validate())
	})
	t.Run("departure_unknown_trucker", func(t *testing.T) {
		order := newValidOrderGraph()
		order.ContainerPoints[0].Departure.ByContactID = "ghost"
		assert.Error(t, order.Validate())
	})
	t.Run("arrival_unknown_trucker", func(t *testing.T) {
		order := newValidOrderGraph()
		order.ContainerPoints[1].Arrival.ByContactID = "ghost"
		assert.Error(t, order.Validate())
	})
	t.Run("duplicate_container_point", func(t *testing.T) {
		order := newValidOrderGraph()
		order.ContainerPoints = append(order.ContainerPoints, &ContainerPoint{
			ContainerID:     "cont1",
			ShippingPointID: "spd",
			ShippingPointBase: ShippingPointBase{
				Status:       ShippingPointStatusPending,
				FreightPoint: FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskLoad}},
			},
		})
		assert.Error(t, order.Validate())
	})
}
