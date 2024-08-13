package facade4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/mocks4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestAddSegmentsTx tests AddSegmentsTx
func TestAddSegmentsTx(t *testing.T) { // TODO: create few test cases
	ctx := context.Background()

	order := mocks4logist.ValidOrderWith3UnassignedContainers(t)

	tx := mocks4logist.MockTx(t)

	request := dto4logist.AddSegmentsRequest{
		Containers: []dto4logist.SegmentContainerData{
			{
				ID: mocks4logist.Container1ID,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks:  []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskLoad},
					ToLoad: &dbo4logist.FreightLoad{NumberOfPallets: 3, GrossWeightKg: 1300, VolumeM3: 2},
				},
			},
			{
				ID: mocks4logist.Container3ID,
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskLoad},
				},
			},
		},
		From: dto4logist.AddSegmentEndpoint{
			AddSegmentParty: dto4logist.AddSegmentParty{
				Counterparty: dbo4logist.SegmentCounterparty{
					ContactID: mocks4logist.Port2dock1ContactID,
					Role:      dbo4logist.CounterpartyRolePickPoint,
				},
			},
		},
		To: dto4logist.AddSegmentEndpoint{
			AddSegmentParty: dto4logist.AddSegmentParty{
				Counterparty: dbo4logist.SegmentCounterparty{
					ContactID: mocks4logist.Dispatcher1warehouse1ContactID,
					Role:      dbo4logist.CounterpartyRoleDispatchPoint,
				},
			},
		},
		By: &dto4logist.AddSegmentParty{
			Counterparty: dbo4logist.SegmentCounterparty{
				ContactID: "trucker1",
				Role:      dbo4logist.CounterpartyRoleTrucker,
			},
		},
	}

	params := &OrderWorkerParams{
		SpaceWorkerParams: &dal4spaceus.SpaceWorkerParams{
			Space: dbo4spaceus.NewSpaceEntry("space1"),
		},
		Order: dbo4logist.NewOrderWithData("space1", "order1", order),
	}

	{ // Pre-checks
		assert.Nil(t, order.Validate())
		assert.Equal(t, 0, len(order.Segments))
		assert.Equal(t, 0, len(order.ShippingPoints))
	}

	if err := addSegmentsTx(ctx, tx, params, request); err != nil {
		t.Fatal("addSegmentsTx() returned unexpected error:", err)
	}

	//b, err := json.MarshalIndent(order, "", "  ")
	//if err != nil {
	//	t.Fatal("json.Marshal() returned unexpected error:", err)
	//}
	//t.Logf("order: %s", string(b))

	order.UpdateCalculatedFields()
	if err := order.Validate(); err != nil {
		t.Error("order is not valid after performing addSegmentsTx():", err)
	}
	assert.True(t, params.Changed.HasChanges())

	const expectedNumberOfSegments = 2 // Because we are adding 2 containers
	assert.Equalf(t, expectedNumberOfSegments, len(order.Segments), "order.Segments: %+v", order.Segments)

	for i, segment := range order.Segments {
		assert.Equal(t, request.Containers[i].ID, segment.ContainerID)
	}

	const expectedNumberOfShippingPoints = 2
	assert.Equalf(t, expectedNumberOfShippingPoints, len(order.ShippingPoints), "order.ShippingPoints: %+v", order.ShippingPoints)
	for _, shippingPoint := range order.ShippingPoints {
		assert.Equal(t, 1, len(shippingPoint.Tasks))
	}
	assert.Equal(t, dbo4logist.ShippingPointTaskPick, order.ShippingPoints[0].Tasks[0])
	assert.Equal(t, dbo4logist.ShippingPointTaskLoad, order.ShippingPoints[1].Tasks[0])
}

func TestAddCounterpartyToOrderIfNeeded(t *testing.T) {
	ctx := context.Background()

	order := mocks4logist.ValidEmptyOrder(t)

	segmentCounterparty := dbo4logist.SegmentCounterparty{
		ContactID: mocks4logist.Dispatcher1warehouse1ContactID,
		Role:      dbo4logist.CounterpartyRoleDispatchPoint,
	}

	tx := mocks4logist.MockTx(t)

	changes, err := addCounterpartyToOrderIfNeeded(ctx, tx, "space1", order, "from", dto4logist.AddSegmentEndpoint{
		AddSegmentParty: dto4logist.AddSegmentParty{
			Counterparty: segmentCounterparty,
		},
	})
	assert.Nil(t, err)
	assert.True(t, changes.Counterparties)

	order.UpdateCalculatedFields()
	err = order.Validate()
	assert.Nil(t, err, "order record is not valid after doing the test")
}
