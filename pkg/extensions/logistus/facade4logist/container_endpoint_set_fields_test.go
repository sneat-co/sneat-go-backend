package facade4logist

import (
	"context"
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/mocks4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestSetContainerEndpointFields(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		return worker(ctx, nil, &OrderWorkerParams{Order: dbo4logist.Order{Dto: &dbo4logist.OrderDbo{}}})
	}

	request := dto4logist.SetContainerEndpointFieldsRequest{
		ContainerPointRequest: dto4logist.ContainerPointRequest{
			OrderRequest:    dto4logist.NewOrderRequest("space1", "order1"),
			ContainerID:     "c1",
			ShippingPointID: "sp1",
		},
	}
	err := SetContainerEndpointFields(nil, request)
	assert.NotNil(t, err) // Fails because container point is missing in empty DTO
}

func Test_txSetContainerEndpointFields(t *testing.T) {
	ctx := context.Background()
	tx := mocks4logist.MockTx(t)

	orderDto := &dbo4logist.OrderDbo{
		WithContainerPoints: dbo4logist.WithContainerPoints{
			ContainerPoints: []*dbo4logist.ContainerPoint{
				{
					ContainerID:     "c1",
					ShippingPointID: "sp1",
					ContainerEndpoints: dbo4logist.ContainerEndpoints{
						Arrival: &dbo4logist.ContainerEndpoint{
							ScheduledDate: "2023-01-01",
						},
						Departure: &dbo4logist.ContainerEndpoint{
							ScheduledDate: "2023-01-02",
						},
					},
				},
			},
		},
	}
	params := &OrderWorkerParams{
		Order: dbo4logist.Order{Dto: orderDto},
	}

	t.Run("success_arrival", func(t *testing.T) {
		request := dto4logist.SetContainerEndpointFieldsRequest{
			ContainerPointRequest: dto4logist.ContainerPointRequest{
				ContainerID:     "c1",
				ShippingPointID: "sp1",
			},
			Side: dbo4logist.EndpointSideArrival,
			Dates: map[string]string{
				"scheduledDate": "2023-01-05",
			},
		}
		err := txSetContainerEndpointFields(ctx, tx, params, request)
		assert.Nil(t, err)
		assert.Equal(t, "2023-01-05", orderDto.ContainerPoints[0].Arrival.ScheduledDate)
		// Departure date should be shifted by 1 day (scheduledDatesDiff was 1 day)
		assert.Equal(t, "2023-01-06", orderDto.ContainerPoints[0].Departure.ScheduledDate)
	})

	t.Run("missing_side", func(t *testing.T) {
		request := dto4logist.SetContainerEndpointFieldsRequest{
			ContainerPointRequest: dto4logist.ContainerPointRequest{
				ContainerID:     "c1",
				ShippingPointID: "sp1",
			},
		}
		err := txSetContainerEndpointFields(ctx, tx, params, request)
		assert.NotNil(t, err)
	})

	t.Run("unknown_side", func(t *testing.T) {
		request := dto4logist.SetContainerEndpointFieldsRequest{
			ContainerPointRequest: dto4logist.ContainerPointRequest{
				ContainerID:     "c1",
				ShippingPointID: "sp1",
			},
			Side: "wrong",
		}
		err := txSetContainerEndpointFields(ctx, tx, params, request)
		assert.NotNil(t, err)
	})

	t.Run("by_contact_id_new", func(t *testing.T) {
		contactID := mocks4logist.Trucker1ContactID
		request := dto4logist.SetContainerEndpointFieldsRequest{
			ContainerPointRequest: dto4logist.ContainerPointRequest{
				ContainerID:     "c1",
				ShippingPointID: "sp1",
			},
			Side:        dbo4logist.EndpointSideArrival,
			ByContactID: &contactID,
		}
		params.SpaceWorkerParams = &dal4spaceus.SpaceWorkerParams{
			Space: dbo4spaceus.NewSpaceEntry("space1"),
		}
		err := txSetContainerEndpointFields(ctx, tx, params, request)
		assert.Nil(t, err)
		assert.Equal(t, contactID, orderDto.ContainerPoints[0].Arrival.ByContactID)
		assert.Equal(t, contactID, orderDto.ContainerPoints[0].Departure.ByContactID)
	})
}
