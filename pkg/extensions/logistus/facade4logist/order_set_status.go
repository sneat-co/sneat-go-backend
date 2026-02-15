package facade4logist

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
)

// SetOrderStatus changes order status
func SetOrderStatus(ctx facade.ContextWithUser, request dto4logist.SetOrderStatusRequest) error {
	return RunOrderWorker(ctx, request.OrderRequest, func(_ facade.ContextWithUser, _ dal.ReadwriteTransaction, params *OrderWorkerParams) error {
		return setOrderStatusTx(params, request)
	})
}

// setOrderStatusTx changes order status in transaction
func setOrderStatusTx(params *OrderWorkerParams, request dto4logist.SetOrderStatusRequest) error {
	if params.Order.Dto.Status == request.Status {
		return nil
	}
	params.Order.Dto.Status = request.Status
	params.Changed.Status = true
	return nil
}
