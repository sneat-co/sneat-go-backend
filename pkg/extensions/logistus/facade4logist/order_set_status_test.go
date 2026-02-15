package facade4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestSetOrderStatus(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		return worker(ctx, nil, &OrderWorkerParams{Order: dbo4logist.Order{Dto: &dbo4logist.OrderDbo{}}})
	}

	request := dto4logist.SetOrderStatusRequest{
		OrderRequest: dto4logist.NewOrderRequest("space1", "order1"),
		Status:       dbo4logist.OrderStatusActive,
	}
	err := SetOrderStatus(nil, request)
	assert.Nil(t, err)
}

func Test_setOrderStatusTx(t *testing.T) {
	orderDto := &dbo4logist.OrderDbo{
		OrderBase: dbo4logist.OrderBase{
			Status: dbo4logist.OrderStatusDraft,
		},
	}
	params := &OrderWorkerParams{
		Order: dbo4logist.Order{Dto: orderDto},
	}

	t.Run("change_status", func(t *testing.T) {
		request := dto4logist.SetOrderStatusRequest{
			Status: dbo4logist.OrderStatusActive,
		}
		err := setOrderStatusTx(params, request)
		assert.Nil(t, err)
		assert.Equal(t, dbo4logist.OrderStatusActive, orderDto.Status)
		assert.True(t, params.Changed.Status)
	})

	t.Run("same_status", func(t *testing.T) {
		params.Changed.Status = false
		request := dto4logist.SetOrderStatusRequest{
			Status: dbo4logist.OrderStatusActive,
		}
		err := setOrderStatusTx(params, request)
		assert.Nil(t, err)
		assert.False(t, params.Changed.Status)
	})
}
