package facade4logist

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// SetContainerPointTask adds/remove task for a container point
func SetContainerPointTask(ctx facade.ContextWithUser, request dto4logist.SetContainerPointTaskRequest) error {
	return RunOrderWorker(ctx, ctx.User(), request.OrderRequest,
		func(_ context.Context, _ dal.ReadwriteTransaction, params *OrderWorkerParams) error {
			return txSetContainerPointTask(params, request)
		},
	)
}

func txSetContainerPointTask(
	params *OrderWorkerParams,
	request dto4logist.SetContainerPointTaskRequest,
) error {
	containerPoint := params.Order.Dto.GetContainerPoint(request.ContainerID, request.ShippingPointID)
	changed := false
	if containerPoint == nil {
		containerPoint = &dbo4logist.ContainerPoint{
			ContainerID:     request.ContainerID,
			ShippingPointID: request.ShippingPointID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
			},
		}
		params.Order.Dto.ContainerPoints = append(params.Order.Dto.ContainerPoints, containerPoint)
		changed = true
	}
	if request.Value {
		if slice.Index(containerPoint.Tasks, request.Task) < 0 {
			containerPoint.Tasks = append(containerPoint.Tasks, request.Task)
			changed = true
		}
	} else {
		if slice.Index(containerPoint.Tasks, request.Task) >= 0 {
			tasks := make([]dbo4logist.ShippingPointTask, 0, len(containerPoint.Tasks)-1)
			for _, task := range containerPoint.Tasks {
				if task != request.Task {
					tasks = append(tasks, task)
				}
			}
			containerPoint.Tasks = tasks
			switch request.Task {
			case dbo4logist.ShippingPointTaskLoad:
				containerPoint.ToLoad = nil
			case dbo4logist.ShippingPointTaskUnload:
				containerPoint.ToUnload = nil
			}
			changed = true
		}
	}
	if changed {
		params.Changed.ContainerPoints = true
	}
	return nil
}
