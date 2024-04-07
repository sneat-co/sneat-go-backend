package facade4logist

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
)

// SetOrderStatus changes order status
func SetOrderStatus(ctx context.Context, user facade.User, request dto4logist.SetOrderStatusRequest) error {
	return RunOrderWorker(ctx, user, request.OrderRequest, func(_ context.Context, _ dal.ReadwriteTransaction, params *OrderWorkerParams) error {
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
