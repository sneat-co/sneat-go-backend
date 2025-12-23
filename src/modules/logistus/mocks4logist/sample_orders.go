package mocks4logist

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-core/coretypes"
	dbmodels2 "github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/with"
)

const (
	Container1ID     = "container1"
	Container1Type   = "10ft"
	Container1Number = "container1number"

	Container2ID     = "container2"
	Container2Type   = "20ft"
	Container2Number = "container2number"

	Container3ID     = "container3"
	Container3Type   = "40ft"
	Container3Number = ""

	ShippingPoint1WithSingleContainerID = "shippingPoint1"
	ShippingPoint2With2ContainersID     = "shippingPoint2"
	ShippingPoint3WithoutContainersID   = "shippingPoint3"

	Port1ContactID      = "port1owner"
	Port1ContactTitle   = "Port 1"
	Port1dock1ContactID = "port1dock1"

	Port2ContactID         = "port2owner"
	Port2ContactTitle      = "Port 2"
	Port2dock1ContactID    = "port2dock1"
	Port2dock2ContactID    = "port2dock2"
	Port2dock2ContactTitle = "Dock 2"

	Trucker1ContactID = "trucker1"

	//Port1dock1shippingPointID = "shippingPointPort1dock1"

	// Port2dock1shippingPointID = shippingPointPort2dock1
	Port2dock1shippingPointID = "shippingPointPort2dock1"
	Port2dock2shippingPointID = "shippingPointPort2dock2"

	Dispatcher1ContactID              = "dispatcher1"
	Dispatcher1ContactTitle           = "WarehouseOperator 1"
	Dispatcher1warehouse1ContactID    = "dispatcher1warehouse1"
	Dispatcher1warehouse1ContactTitle = "WarehouseOperator 1"
	Dispatcher2ContactID              = "dispatcher2"
	Dispatcher2ContactTitle           = "WarehouseOperator 2"
	Dispatcher2warehouse1ContactID    = "dispatcher2warehouse1"
	Dispatcher2warehouse1ContactTitle = "WarehouseOperator 1"
	Dispatcher2warehouse2ContactID    = "dispatcher2warehouse2"
	Dispatcher2warehouse2ContactTitle = "WarehouseOperator 2"
)

const (
	Dto1ContainerPointsCount = 5
	Dto1ShippingPointsCount  = 5
)

// ValidOrderDto1 returns a valid order with 3 containers, 3 shipping points, 4 container points and 1 port
func ValidOrderDto1(t *testing.T) (order *dbo4logist.OrderDbo) {
	order = ValidEmptyOrder(t)

	order.Contacts = append(order.Contacts, []*dbo4logist.OrderContact{
		{
			ID:        Port2ContactID,
			Type:      briefs4contactus.ContactTypeCompany,
			Title:     "Port 2 owner",
			CountryID: "IE",
		},
		{
			ID:        Port2dock1ContactID,
			ParentID:  Port2ContactID,
			Type:      briefs4contactus.ContactTypeLocation,
			Title:     "Port 2 dock 1",
			CountryID: "IE",
		},
		{
			ID:        Port2dock2ContactID,
			Type:      briefs4contactus.ContactTypeLocation,
			ParentID:  Port2ContactID,
			Title:     "Port 2 dock 2",
			CountryID: "IE",
		},
		{
			ID:        Port1ContactID,
			Type:      briefs4contactus.ContactTypeCompany,
			Title:     "Port 1 owner",
			CountryID: "IE",
		},
		{
			ID:        Port1dock1ContactID,
			ParentID:  Port1ContactID,
			Type:      briefs4contactus.ContactTypeLocation,
			Title:     "Port 1 dock 1",
			CountryID: "IE",
		},
		{
			ID:        Dispatcher1ContactID,
			Type:      briefs4contactus.ContactTypeCompany,
			Title:     Dispatcher1ContactTitle,
			CountryID: "IE",
		},
		{
			ID:        Dispatcher1warehouse1ContactID,
			ParentID:  Dispatcher1ContactID,
			Type:      briefs4contactus.ContactTypeLocation,
			Title:     Dispatcher1warehouse1ContactTitle,
			CountryID: "IE",
		},
		{
			ID:        Dispatcher2ContactID,
			Type:      briefs4contactus.ContactTypeCompany,
			Title:     Dispatcher2ContactTitle,
			CountryID: "IE",
		},
		{
			ID:        Dispatcher2warehouse1ContactID,
			ParentID:  Dispatcher2ContactID,
			Type:      briefs4contactus.ContactTypeLocation,
			Title:     Dispatcher2warehouse1ContactTitle,
			CountryID: "IE",
		},
		{
			ID:        Dispatcher2warehouse2ContactID,
			ParentID:  Dispatcher2ContactID,
			Type:      briefs4contactus.ContactTypeLocation,
			Title:     Dispatcher2warehouse2ContactTitle,
			CountryID: "IE",
		},
	}...)

	order.Counterparties = append(order.Counterparties, []*dbo4logist.OrderCounterparty{
		{ContactID: Port2ContactID, Role: dbo4logist.CounterpartyRolePortFrom},
		{ContactID: Port2dock1ContactID, Role: dbo4logist.CounterpartyRolePickPoint},
		{ContactID: Port2dock1ContactID, Role: dbo4logist.CounterpartyRoleDropPoint},
		{ContactID: Port2dock2ContactID, Role: dbo4logist.CounterpartyRolePickPoint},
		{ContactID: Port2dock2ContactID, Role: dbo4logist.CounterpartyRoleDropPoint},
		{ContactID: Port1ContactID, Role: dbo4logist.CounterpartyRolePortTo},
		{ContactID: Port1dock1ContactID, Role: dbo4logist.CounterpartyRoleDropPoint},
		{ContactID: Dispatcher1ContactID, Role: dbo4logist.CounterpartyRoleDispatcher},
		{ContactID: Dispatcher1warehouse1ContactID, Role: dbo4logist.CounterpartyRoleDispatchPoint},
		{ContactID: Dispatcher2ContactID, Role: dbo4logist.CounterpartyRoleDispatcher},
		{ContactID: Dispatcher2warehouse1ContactID, Role: dbo4logist.CounterpartyRoleDispatchPoint},
		{ContactID: Dispatcher2warehouse2ContactID, Role: dbo4logist.CounterpartyRoleDispatchPoint},
	}...)
	fixCounterpartiesFromContacts(order)
	order.Containers = []*dbo4logist.OrderContainer{
		{
			ID: Container1ID,
			OrderContainerBase: dbo4logist.OrderContainerBase{
				Type:   Container1Type,
				Number: Container1Number,
			},
		},
		{
			ID: Container2ID,
			OrderContainerBase: dbo4logist.OrderContainerBase{
				Type:   Container2Type,
				Number: Container2Number,
			},
		},
		{
			ID: Container3ID,
			OrderContainerBase: dbo4logist.OrderContainerBase{
				Type:   Container3Type,
				Number: Container3Number,
			},
		},
	}
	order.Segments = []*dbo4logist.ContainerSegment{
		{
			Dates: &dbo4logist.SegmentDates{
				Arrives: "2020-12-31",
			},
			ContainerSegmentKey: dbo4logist.ContainerSegmentKey{
				ContainerID: Container2ID,
				From: dbo4logist.SegmentEndpoint{
					ShippingPointID: Port2dock2shippingPointID,
					SegmentCounterparty: dbo4logist.SegmentCounterparty{
						ContactID: Port2dock1ContactID,
						Role:      dbo4logist.CounterpartyRolePickPoint,
					},
				},
				To: dbo4logist.SegmentEndpoint{
					ShippingPointID: ShippingPoint1WithSingleContainerID,
					SegmentCounterparty: dbo4logist.SegmentCounterparty{
						ContactID: Dispatcher1warehouse1ContactID,
						Role:      dbo4logist.CounterpartyRoleDispatchPoint,
					},
				},
			},
		},
	}
	order.ContainerPoints = []*dbo4logist.ContainerPoint{
		{
			ContainerID:     Container1ID,
			ShippingPointID: Port2dock1shippingPointID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskPick},
				},
			},
		},
		{
			ContainerID:     Container2ID,
			ShippingPointID: Port2dock2shippingPointID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskPick},
				},
			},
		},
		{
			ContainerID:     Container2ID,
			ShippingPointID: ShippingPoint1WithSingleContainerID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskLoad},
				},
			},
			ContainerEndpoints: dbo4logist.ContainerEndpoints{
				Arrival: &dbo4logist.ContainerEndpoint{
					ScheduledDate: "2020-12-31",
				},
			},
		},
		{
			ContainerID:     Container2ID,
			ShippingPointID: ShippingPoint2With2ContainersID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskLoad, dbo4logist.ShippingPointTaskUnload},
				},
			},
		},
		{
			ContainerID:     Container1ID,
			ShippingPointID: ShippingPoint2With2ContainersID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskLoad},
				},
			},
		},
	}
	order.ShippingPoints = []*dbo4logist.OrderShippingPoint{
		{
			ID: Port2dock1shippingPointID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusCompleted,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []string{dbo4logist.ShippingPointTaskPick, dbo4logist.ShippingPointTaskDrop},
				},
			},
			Counterparty: dbo4logist.ShippingPointCounterparty{
				ContactID: Port2ContactID,
				Title:     Port2ContactTitle,
			},
			Location: &dbo4logist.ShippingPointLocation{
				ContactID: Port2dock1ContactID,
				Title:     "Port 2 dock 1",
				Address: &dbmodels2.Address{
					CountryID: "IE",
					Lines:     "Dock 1\nPort 2",
				},
			},
		},
		{
			ID: Port2dock2shippingPointID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusCompleted,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []string{dbo4logist.ShippingPointTaskPick, dbo4logist.ShippingPointTaskDrop},
				},
			},
			Counterparty: dbo4logist.ShippingPointCounterparty{
				ContactID: Port2ContactID,
				Title:     Port2ContactTitle,
			},
			Location: &dbo4logist.ShippingPointLocation{
				ContactID: Port2dock2ContactID,
				Title:     Port2dock2ContactTitle,
				Address: &dbmodels2.Address{
					CountryID: "IE",
					Lines:     "Dock 2\nPort 2",
				},
			},
		},
		{
			ID: ShippingPoint1WithSingleContainerID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []string{dbo4logist.ShippingPointTaskLoad},
				},
			},
			Counterparty: dbo4logist.ShippingPointCounterparty{
				ContactID: Dispatcher1ContactID,
				Title:     "WarehouseOperator 1",
			},
		},
		{
			ID: ShippingPoint2With2ContainersID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []string{dbo4logist.ShippingPointTaskLoad},
				},
			},
			Counterparty: dbo4logist.ShippingPointCounterparty{
				ContactID: Dispatcher2warehouse1ContactID,
				Title:     "WarehouseOperator 2",
			},
		},
		{
			ID: ShippingPoint3WithoutContainersID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []string{dbo4logist.ShippingPointTaskLoad},
				},
			},
			Counterparty: dbo4logist.ShippingPointCounterparty{
				ContactID: Dispatcher2ContactID,
				Title:     Dispatcher2ContactTitle,
			},
		},
	}
	order.UpdateCalculatedFields()
	if err := order.Validate(); err != nil {
		t.Errorf("this test order must be valid: %v", err)
	}
	assert.Equal(t, 13, len(order.Counterparties))
	assert.Equal(t, Dto1ShippingPointsCount, len(order.ShippingPoints))
	assert.Equal(t, 1, len(order.GetCounterpartiesByRole(dbo4logist.CounterpartyRolePortFrom)))
	assert.Equal(t, 1, len(order.GetCounterpartiesByRole(dbo4logist.CounterpartyRolePortTo)))
	assert.Equal(t, 1, len(order.GetCounterpartiesByRole(dbo4logist.CounterpartyRolePortTo)))
	return order
}

func fixCounterpartiesFromContacts(order *dbo4logist.OrderDbo) {
	for _, cp := range order.Counterparties {
		contact := order.MustGetContactByID(cp.ContactID)
		cp.CountryID = contact.CountryID
		cp.Title = contact.Title
		if contact.ParentID != "" {
			_, parent := order.GetCounterpartyByContactID(contact.ParentID)
			cp.Parent = &dbo4logist.CounterpartyParent{
				ContactID: parent.ContactID,
				Role:      parent.Role,
			}
		}
	}
}

// ValidOrderWith3UnassignedContainers returns a valid order with 3 unassigned containers
func ValidOrderWith3UnassignedContainers(t *testing.T) (order *dbo4logist.OrderDbo) {
	order = ValidEmptyOrder(t)
	order.Containers = []*dbo4logist.OrderContainer{
		{
			ID: Container1ID,
			OrderContainerBase: dbo4logist.OrderContainerBase{
				Type:   Container1Type,
				Number: "C1",
			},
		},
		{
			ID: Container2ID,
			OrderContainerBase: dbo4logist.OrderContainerBase{
				Type:   Container2Type,
				Number: "C2",
			},
		},
		{
			ID: Container3ID,
			OrderContainerBase: dbo4logist.OrderContainerBase{
				Type:   Container3Type,
				Number: "C4",
			},
		},
	}
	if err := order.Validate(); err != nil {
		t.Errorf("this test order must be valid: %v", err)
	}
	return order
}

// ValidEmptyOrder returns a valid empty order
func ValidEmptyOrder(t *testing.T) (order *dbo4logist.OrderDbo) {
	modified := dbmodels2.Modified{
		By: "unit-test",
		At: time.Now(),
	}

	order = &dbo4logist.OrderDbo{
		WithModified: dbmodels2.WithModified{
			CreatedFields: with.CreatedFields{
				CreatedAtField: with.CreatedAtField{
					CreatedAt: modified.At,
				},
				CreatedByField: with.CreatedByField{
					CreatedBy: modified.By,
				},
			},
			UpdatedFields: with.UpdatedFields{
				UpdatedAt: modified.At,
				UpdatedBy: modified.By,
			},
		},
		WithSpaceID: dbmodels2.WithSpaceID{
			SpaceID: "space-1",
		},
		WithSpaceIDs: dbmodels2.WithSpaceIDs{
			SpaceIDs: []coretypes.SpaceID{"space-1", "space-2"},
		},
		WithUserIDs: dbmodels2.WithUserIDs{
			UserIDs: []string{"user-1", "user-2"},
		},
		WithOrderContacts: dbo4logist.WithOrderContacts{
			Contacts: []*dbo4logist.OrderContact{
				{
					ID:        "buyer1",
					Type:      "company",
					Title:     "Buyer 1",
					CountryID: "ES",
				},
			},
		},
		OrderBase: dbo4logist.OrderBase{
			Status:    "active",
			Direction: "export",
			WithCounterparties: dbo4logist.WithCounterparties{
				Counterparties: []*dbo4logist.OrderCounterparty{
					{
						ContactID: "buyer1",
						Role:      "buyer",
					},
				},
			},
		},
	}
	fixCounterpartiesFromContacts(order)
	order.UpdateCalculatedFields()
	if err := order.Validate(); err != nil {
		t.Errorf("this test order must be valid: %v", err)
	}
	return order
}
