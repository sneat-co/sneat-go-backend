package facade4logist

import (
	"context"
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/mocks4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestUpdateShippingPoint(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		return worker(ctx, nil, &OrderWorkerParams{Order: dbo4logist.Order{Dto: &dbo4logist.OrderDbo{}}})
	}

	request := dto4logist.UpdateShippingPointRequest{
		OrderShippingPointRequest: dto4logist.OrderShippingPointRequest{
			OrderRequest:    dto4logist.NewOrderRequest("space1", "order1"),
			ShippingPointID: "sp1",
		},
	}
	err := UpdateShippingPoint(nil, request)
	assert.Nil(t, err)
}

func Test_txUpdateShippingPoint(t *testing.T) {
	ctx := context.Background()
	tx := mocks4logist.MockTx(t)
	orderDto := &dbo4logist.OrderDbo{
		WithShippingPoints: dbo4logist.WithShippingPoints{
			ShippingPoints: []*dbo4logist.OrderShippingPoint{
				{ID: "sp1"},
			},
		},
	}
	params := &OrderWorkerParams{
		Order: dbo4logist.Order{Dto: orderDto},
	}

	t.Run("success", func(t *testing.T) {
		request := dto4logist.UpdateShippingPointRequest{
			OrderShippingPointRequest: dto4logist.OrderShippingPointRequest{
				ShippingPointID: "sp1",
			},
			SetFieldsRequest: dto4logist.SetFieldsRequest{
				SetDates: map[string]string{
					"scheduledStartDate": "2023-01-01",
					"scheduledEndDate":   "2023-01-02",
				},
				SetStrings: map[string]string{
					"notes": "some notes",
				},
			},
		}
		err := txUpdateShippingPoint(ctx, tx, params, request)
		assert.Nil(t, err)
		assert.Equal(t, "2023-01-01", orderDto.ShippingPoints[0].ScheduledStartDate)
		assert.Equal(t, "2023-01-02", orderDto.ShippingPoints[0].ScheduledEndDate)
		assert.Equal(t, "some notes", orderDto.ShippingPoints[0].Notes)
		assert.True(t, params.Changed.ShippingPoints)
	})

	t.Run("unknown_date_field", func(t *testing.T) {
		request := dto4logist.UpdateShippingPointRequest{
			OrderShippingPointRequest: dto4logist.OrderShippingPointRequest{
				ShippingPointID: "sp1",
			},
			SetFieldsRequest: dto4logist.SetFieldsRequest{
				SetDates: map[string]string{
					"wrong": "2023-01-01",
				},
			},
		}
		err := txUpdateShippingPoint(ctx, tx, params, request)
		assert.NotNil(t, err)
	})

	t.Run("unknown_string_field", func(t *testing.T) {
		request := dto4logist.UpdateShippingPointRequest{
			OrderShippingPointRequest: dto4logist.OrderShippingPointRequest{
				ShippingPointID: "sp1",
			},
			SetFieldsRequest: dto4logist.SetFieldsRequest{
				SetStrings: map[string]string{
					"wrong": "value",
				},
			},
		}
		err := txUpdateShippingPoint(ctx, tx, params, request)
		assert.NotNil(t, err)
	})
}
