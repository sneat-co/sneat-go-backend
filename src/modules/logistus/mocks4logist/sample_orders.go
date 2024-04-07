package mocks4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	dbmodels2 "github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/with"
	"testing"
	"time"
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

	Port1dock1shippingPointID = "shippingPointPort1dock1"
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
func ValidOrderDto1(t *testing.T) (order *models4logist.OrderDto) {
	order = ValidEmptyOrder(t)

	order.Contacts = append(order.Contacts, []*models4logist.OrderContact{
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

	order.Counterparties = append(order.Counterparties, []*models4logist.OrderCounterparty{
		{ContactID: Port2ContactID, Role: models4logist.CounterpartyRolePortFrom},
		{ContactID: Port2dock1ContactID, Role: models4logist.CounterpartyRolePickPoint},
		{ContactID: Port2dock1ContactID, Role: models4logist.CounterpartyRoleDropPoint},
		{ContactID: Port2dock2ContactID, Role: models4logist.CounterpartyRolePickPoint},
		{ContactID: Port2dock2ContactID, Role: models4logist.CounterpartyRoleDropPoint},
		{ContactID: Port1ContactID, Role: models4logist.CounterpartyRolePortTo},
		{ContactID: Port1dock1ContactID, Role: models4logist.CounterpartyRoleDropPoint},
		{ContactID: Dispatcher1ContactID, Role: models4logist.CounterpartyRoleDispatcher},
		{ContactID: Dispatcher1warehouse1ContactID, Role: models4logist.CounterpartyRoleDispatchPoint},
		{ContactID: Dispatcher2ContactID, Role: models4logist.CounterpartyRoleDispatcher},
		{ContactID: Dispatcher2warehouse1ContactID, Role: models4logist.CounterpartyRoleDispatchPoint},
		{ContactID: Dispatcher2warehouse2ContactID, Role: models4logist.CounterpartyRoleDispatchPoint},
	}...)
	fixCounterpartiesFromContacts(order)
	order.Containers = []*models4logist.OrderContainer{
		{
			ID: Container1ID,
			OrderContainerBase: models4logist.OrderContainerBase{
				Type:   Container1Type,
				Number: Container1Number,
			},
		},
		{
			ID: Container2ID,
			OrderContainerBase: models4logist.OrderContainerBase{
				Type:   Container2Type,
				Number: Container2Number,
			},
		},
		{
			ID: Container3ID,
			OrderContainerBase: models4logist.OrderContainerBase{
				Type:   Container3Type,
				Number: Container3Number,
			},
		},
	}
	order.Segments = []*models4logist.ContainerSegment{
		{
			Dates: &models4logist.SegmentDates{
				Arrives: "2020-12-31",
			},
			ContainerSegmentKey: models4logist.ContainerSegmentKey{
				ContainerID: Container2ID,
				From: models4logist.SegmentEndpoint{
					ShippingPointID: Port2dock2shippingPointID,
					SegmentCounterparty: models4logist.SegmentCounterparty{
						ContactID: Port2dock1ContactID,
						Role:      models4logist.CounterpartyRolePickPoint,
					},
				},
				To: models4logist.SegmentEndpoint{
					ShippingPointID: ShippingPoint1WithSingleContainerID,
					SegmentCounterparty: models4logist.SegmentCounterparty{
						ContactID: Dispatcher1warehouse1ContactID,
						Role:      models4logist.CounterpartyRoleDispatchPoint,
					},
				},
			},
		},
	}
	order.ContainerPoints = []*models4logist.ContainerPoint{
		{
			ContainerID:     Container1ID,
			ShippingPointID: Port2dock1shippingPointID,
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusPending,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []models4logist.ShippingPointTask{models4logist.ShippingPointTaskPick},
				},
			},
		},
		{
			ContainerID:     Container2ID,
			ShippingPointID: Port2dock2shippingPointID,
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusPending,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []models4logist.ShippingPointTask{models4logist.ShippingPointTaskPick},
				},
			},
		},
		{
			ContainerID:     Container2ID,
			ShippingPointID: ShippingPoint1WithSingleContainerID,
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusPending,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []models4logist.ShippingPointTask{models4logist.ShippingPointTaskLoad},
				},
			},
			ContainerEndpoints: models4logist.ContainerEndpoints{
				Arrival: &models4logist.ContainerEndpoint{
					ScheduledDate: "2020-12-31",
				},
			},
		},
		{
			ContainerID:     Container2ID,
			ShippingPointID: ShippingPoint2With2ContainersID,
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusPending,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []models4logist.ShippingPointTask{models4logist.ShippingPointTaskLoad, models4logist.ShippingPointTaskUnload},
				},
			},
		},
		{
			ContainerID:     Container1ID,
			ShippingPointID: ShippingPoint2With2ContainersID,
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusPending,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []models4logist.ShippingPointTask{models4logist.ShippingPointTaskLoad},
				},
			},
		},
	}
	order.ShippingPoints = []*models4logist.OrderShippingPoint{
		{
			ID: Port2dock1shippingPointID,
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusCompleted,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []string{models4logist.ShippingPointTaskPick, models4logist.ShippingPointTaskDrop},
				},
			},
			Counterparty: models4logist.ShippingPointCounterparty{
				ContactID: Port2ContactID,
				Title:     Port2ContactTitle,
			},
			Location: &models4logist.ShippingPointLocation{
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
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusCompleted,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []string{models4logist.ShippingPointTaskPick, models4logist.ShippingPointTaskDrop},
				},
			},
			Counterparty: models4logist.ShippingPointCounterparty{
				ContactID: Port2ContactID,
				Title:     Port2ContactTitle,
			},
			Location: &models4logist.ShippingPointLocation{
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
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusPending,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []string{models4logist.ShippingPointTaskLoad},
				},
			},
			Counterparty: models4logist.ShippingPointCounterparty{
				ContactID: Dispatcher1ContactID,
				Title:     "WarehouseOperator 1",
			},
		},
		{
			ID: ShippingPoint2With2ContainersID,
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusPending,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []string{models4logist.ShippingPointTaskLoad},
				},
			},
			Counterparty: models4logist.ShippingPointCounterparty{
				ContactID: Dispatcher2warehouse1ContactID,
				Title:     "WarehouseOperator 2",
			},
		},
		{
			ID: ShippingPoint3WithoutContainersID,
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: models4logist.ShippingPointStatusPending,
				FreightPoint: models4logist.FreightPoint{
					Tasks: []string{models4logist.ShippingPointTaskLoad},
				},
			},
			Counterparty: models4logist.ShippingPointCounterparty{
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
	assert.Equal(t, 1, len(order.GetCounterpartiesByRole(models4logist.CounterpartyRolePortFrom)))
	assert.Equal(t, 1, len(order.GetCounterpartiesByRole(models4logist.CounterpartyRolePortTo)))
	assert.Equal(t, 1, len(order.GetCounterpartiesByRole(models4logist.CounterpartyRolePortTo)))
	return order
}

func fixCounterpartiesFromContacts(order *models4logist.OrderDto) {
	for _, cp := range order.Counterparties {
		contact := order.MustGetContactByID(cp.ContactID)
		cp.CountryID = contact.CountryID
		cp.Title = contact.Title
		if contact.ParentID != "" {
			_, parent := order.WithCounterparties.GetCounterpartyByContactID(contact.ParentID)
			cp.Parent = &models4logist.CounterpartyParent{
				ContactID: parent.ContactID,
				Role:      parent.Role,
			}
		}
	}
}

// ValidOrderWith3UnassignedContainers returns a valid order with 3 unassigned containers
func ValidOrderWith3UnassignedContainers(t *testing.T) (order *models4logist.OrderDto) {
	order = ValidEmptyOrder(t)
	order.Containers = []*models4logist.OrderContainer{
		{
			ID: Container1ID,
			OrderContainerBase: models4logist.OrderContainerBase{
				Type:   Container1Type,
				Number: "C1",
			},
		},
		{
			ID: Container2ID,
			OrderContainerBase: models4logist.OrderContainerBase{
				Type:   Container2Type,
				Number: "C2",
			},
		},
		{
			ID: Container3ID,
			OrderContainerBase: models4logist.OrderContainerBase{
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
func ValidEmptyOrder(t *testing.T) (order *models4logist.OrderDto) {
	modified := dbmodels2.Modified{
		By: "unit-test",
		At: time.Now(),
	}

	order = &models4logist.OrderDto{
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
		WithTeamID: dbmodels2.WithTeamID{
			TeamID: "team-1",
		},
		WithTeamIDs: dbmodels2.WithTeamIDs{
			TeamIDs: []string{"team-1", "team-2"},
		},
		WithUserIDs: dbmodels2.WithUserIDs{
			UserIDs: []string{"user-1", "user-2"},
		},
		WithOrderContacts: models4logist.WithOrderContacts{
			Contacts: []*models4logist.OrderContact{
				{
					ID:        "buyer1",
					Type:      "company",
					Title:     "Buyer 1",
					CountryID: "ES",
				},
			},
		},
		OrderBase: models4logist.OrderBase{
			Status:    "active",
			Direction: "export",
			WithCounterparties: models4logist.WithCounterparties{
				Counterparties: []*models4logist.OrderCounterparty{
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
