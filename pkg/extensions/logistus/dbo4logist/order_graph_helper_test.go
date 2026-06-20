package dbo4logist

import (
	"time"

	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// newValidOrderGraph builds a fully wired, valid OrderDbo graph with
// counterparties, contacts, containers, shipping points, container points
// and segments all cross-referencing each other by IDs.
//
// Graph shape:
//   - dispatcher company "disp" + dispatch_point location "dp" (parent=disp)
//   - receiver company "recv" + receive_point location "rp" (parent=recv)
//   - trucker "truck"
//   - one container "cont1"
//   - two shipping points: "spd" (load @ dp) and "spr" (unload @ rp)
//   - two container points linking cont1 to both shipping points
//   - one segment cont1: dp -> rp
func newValidOrderGraph() *OrderDbo {
	order := &OrderDbo{
		WithModified: dbmodels.NewWithModified(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "u1"),
		WithSpaceID:  dbmodels.WithSpaceID{SpaceID: coretypes.SpaceID("s1")},
		WithSpaceIDs: dbmodels.WithSpaceIDs{SpaceIDs: []coretypes.SpaceID{coretypes.SpaceID("s1")}},
		WithUserIDs:  dbmodels.WithUserIDs{UserIDs: []string{"u1"}},
		OrderBase: OrderBase{
			Status:    OrderStatusActive,
			Direction: OrderDirectionImport,
			WithCounterparties: WithCounterparties{
				Counterparties: []*OrderCounterparty{
					{ContactID: "disp", Role: CounterpartyRoleDispatcher, CountryID: "US", Title: "Dispatcher Co"},
					{ContactID: "dp", Role: CounterpartyRoleDispatchPoint, CountryID: "US", Title: "Dispatch Point",
						Parent: &CounterpartyParent{ContactID: "disp", Role: CounterpartyRoleDispatcher}},
					{ContactID: "recv", Role: CounterpartyRoleReceiver, CountryID: "GB", Title: "Receiver Co"},
					{ContactID: "rp", Role: CounterpartyRoleReceivePoint, CountryID: "GB", Title: "Receive Point",
						Parent: &CounterpartyParent{ContactID: "recv", Role: CounterpartyRoleReceiver}},
					{ContactID: "truck", Role: CounterpartyRoleTrucker, CountryID: "US", Title: "Trucker"},
				},
			},
		},
		WithOrderContacts: WithOrderContacts{
			Contacts: []*OrderContact{
				{ID: "disp", Type: briefs4contactus.ContactTypeCompany, Title: "Dispatcher Co", CountryID: "US"},
				{ID: "dp", Type: briefs4contactus.ContactTypeLocation, Title: "Dispatch Point", CountryID: "US", ParentID: "disp"},
				{ID: "recv", Type: briefs4contactus.ContactTypeCompany, Title: "Receiver Co", CountryID: "GB"},
				{ID: "rp", Type: briefs4contactus.ContactTypeLocation, Title: "Receive Point", CountryID: "GB", ParentID: "recv"},
				{ID: "truck", Type: briefs4contactus.ContactTypeCompany, Title: "Trucker", CountryID: "US"},
			},
		},
		WithOrderContainers: WithOrderContainers{
			Containers: []*OrderContainer{
				{ID: "cont1", OrderContainerBase: OrderContainerBase{Type: ContainerType20ft}},
			},
		},
		WithShippingPoints: WithShippingPoints{
			ShippingPoints: []*OrderShippingPoint{
				{
					ID: "spd",
					ShippingPointBase: ShippingPointBase{
						Status:       ShippingPointStatusPending,
						FreightPoint: FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskLoad}},
					},
					ScheduledStartDate: "2023-01-01",
					Counterparty:       ShippingPointCounterparty{ContactID: "disp", Title: "Dispatcher Co"},
					Location:           &ShippingPointLocation{ContactID: "dp", Title: "Dispatch Point", Address: &dbmodels.Address{CountryID: "US"}},
				},
				{
					ID: "spr",
					ShippingPointBase: ShippingPointBase{
						Status:       ShippingPointStatusPending,
						FreightPoint: FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskUnload}},
					},
					ScheduledEndDate: "2023-01-10",
					Counterparty:     ShippingPointCounterparty{ContactID: "recv", Title: "Receiver Co"},
					Location:         &ShippingPointLocation{ContactID: "rp", Title: "Receive Point", Address: &dbmodels.Address{CountryID: "GB"}},
				},
			},
		},
		WithContainerPoints: WithContainerPoints{
			ContainerPoints: []*ContainerPoint{
				{
					ContainerID:     "cont1",
					ShippingPointID: "spd",
					ShippingPointBase: ShippingPointBase{
						Status:       ShippingPointStatusPending,
						FreightPoint: FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskLoad}},
					},
					ContainerEndpoints: ContainerEndpoints{
						Departure: &ContainerEndpoint{ScheduledDate: "2023-01-01"},
					},
				},
				{
					ContainerID:     "cont1",
					ShippingPointID: "spr",
					ShippingPointBase: ShippingPointBase{
						Status:       ShippingPointStatusPending,
						FreightPoint: FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskUnload}},
					},
					ContainerEndpoints: ContainerEndpoints{
						Arrival: &ContainerEndpoint{ScheduledDate: "2023-01-10"},
					},
				},
			},
		},
		WithSegments: WithSegments{
			Segments: []*ContainerSegment{
				{
					ContainerSegmentKey: ContainerSegmentKey{
						ContainerID: "cont1",
						From: SegmentEndpoint{
							SegmentCounterparty: SegmentCounterparty{ContactID: "dp", Role: CounterpartyRoleDispatchPoint},
							ShippingPointID:     "spd",
						},
						To: SegmentEndpoint{
							SegmentCounterparty: SegmentCounterparty{ContactID: "rp", Role: CounterpartyRoleReceivePoint},
							ShippingPointID:     "spr",
						},
					},
					Dates: &SegmentDates{Departs: "2023-01-01", Arrives: "2023-01-10"},
				},
			},
		},
	}
	order.CountryIDs = []string{"US", "GB"}
	order.UpdateKeys()
	return order
}
