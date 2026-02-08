package facade4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestSetContainerFields(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		return worker(ctx, nil, &OrderWorkerParams{Order: dbo4logist.Order{Dto: &dbo4logist.OrderDbo{}}})
	}

	request := dto4logist.SetContainerFieldsRequest{
		ContainerRequest: dto4logist.ContainerRequest{
			OrderRequest: dto4logist.NewOrderRequest("space1", "order1"),
			ContainerID:  "c1",
		},
	}
	err := SetContainerFields(nil, request)
	assert.NotNil(t, err)
}

func Test_txSetContainerFields(t *testing.T) {
	orderDto := &dbo4logist.OrderDbo{
		WithOrderContainers: dbo4logist.WithOrderContainers{
			Containers: []*dbo4logist.OrderContainer{
				{ID: "c1"},
			},
		},
	}
	params := &OrderWorkerParams{
		Order: dbo4logist.Order{Dto: orderDto},
	}

	t.Run("success", func(t *testing.T) {
		request := dto4logist.SetContainerFieldsRequest{
			ContainerRequest: dto4logist.ContainerRequest{
				ContainerID: "c1",
			},
			SetFieldsRequest: dto4logist.SetFieldsRequest{
				SetStrings: map[string]string{
					"number":       "C123",
					"instructions": "Handle with care",
				},
			},
		}
		err := txSetContainerFields(params, request)
		assert.Nil(t, err)
		assert.Equal(t, "C123", orderDto.Containers[0].Number)
		assert.Equal(t, "Handle with care", orderDto.Containers[0].Instructions)
		assert.True(t, params.Changed.Containers)
	})

	t.Run("container_not_found", func(t *testing.T) {
		request := dto4logist.SetContainerFieldsRequest{
			ContainerRequest: dto4logist.ContainerRequest{
				ContainerID: "c2",
			},
		}
		err := txSetContainerFields(params, request)
		assert.NotNil(t, err)
	})

	t.Run("unknown_field", func(t *testing.T) {
		request := dto4logist.SetContainerFieldsRequest{
			ContainerRequest: dto4logist.ContainerRequest{
				ContainerID: "c1",
			},
			SetFieldsRequest: dto4logist.SetFieldsRequest{
				SetStrings: map[string]string{
					"wrong": "value",
				},
			},
		}
		err := txSetContainerFields(params, request)
		assert.NotNil(t, err)
	})
}
