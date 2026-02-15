package facade4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestSetContainerPointFields(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		return worker(ctx, nil, &OrderWorkerParams{Order: dbo4logist.Order{Dto: &dbo4logist.OrderDbo{}}})
	}

	request := dto4logist.SetContainerPointFieldsRequest{
		ContainerPointRequest: dto4logist.ContainerPointRequest{
			OrderRequest:    dto4logist.NewOrderRequest("space1", "order1"),
			ContainerID:     "c1",
			ShippingPointID: "sp1",
		},
	}
	err := SetContainerPointFields(nil, request)
	assert.Nil(t, err)
}

func Test_txSetContainerPointFields(t *testing.T) {
	orderDto := &dbo4logist.OrderDbo{
		WithContainerPoints: dbo4logist.WithContainerPoints{
			ContainerPoints: []*dbo4logist.ContainerPoint{
				{
					ContainerID:     "c1",
					ShippingPointID: "sp1",
					ShippingPointBase: dbo4logist.ShippingPointBase{
						Notes: "old notes",
					},
					RefNumber: "old ref",
				},
			},
		},
	}
	params := &OrderWorkerParams{
		Order: dbo4logist.Order{Dto: orderDto},
	}

	t.Run("success_existing", func(t *testing.T) {
		request := dto4logist.SetContainerPointFieldsRequest{
			ContainerPointRequest: dto4logist.ContainerPointRequest{
				ContainerID:     "c1",
				ShippingPointID: "sp1",
			},
			SetStrings: map[string]string{
				"notes":     "new notes",
				"refNumber": "new ref",
			},
		}
		err := txSetContainerPointFields(params, request)
		assert.Nil(t, err)
		assert.Equal(t, "new notes", orderDto.ContainerPoints[0].Notes)
		assert.Equal(t, "new ref", orderDto.ContainerPoints[0].RefNumber)
		assert.True(t, params.Changed.ContainerPoints)
	})

	t.Run("unknown_field", func(t *testing.T) {
		request := dto4logist.SetContainerPointFieldsRequest{
			ContainerPointRequest: dto4logist.ContainerPointRequest{
				ContainerID:     "c1",
				ShippingPointID: "sp1",
			},
			SetStrings: map[string]string{
				"wrong": "value",
			},
		}
		err := txSetContainerPointFields(params, request)
		assert.NotNil(t, err)
	})
}
