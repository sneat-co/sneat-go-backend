package facade4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestDeleteContainer(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		return worker(ctx, nil, &OrderWorkerParams{Order: dbo4logist.Order{Dto: &dbo4logist.OrderDbo{}}})
	}

	request := dto4logist.ContainerRequest{
		OrderRequest: dto4logist.NewOrderRequest("space1", "order1"),
		ContainerID:  "c1",
	}
	err := DeleteContainer(nil, request)
	assert.Nil(t, err)
}

func Test_deleteContainerTx(t *testing.T) {
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
	request := dto4logist.ContainerRequest{
		ContainerID: "c1",
	}
	err := deleteContainerTx(request, params)
	assert.Nil(t, err)
	assert.True(t, params.Changed.Containers)
	assert.True(t, params.Changed.ContainerPoints)
	assert.Empty(t, orderDto.Containers)

	// Test non-existing container
	params.Changed.Containers = false
	err = deleteContainerTx(dto4logist.ContainerRequest{ContainerID: "c2"}, params)
	assert.Nil(t, err)
	assert.False(t, params.Changed.Containers)
}
