package facade4logist

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
)

// AddContainerPoints adds container point to an order
func AddContainerPoints(ctx facade.ContextWithUser, request dto4logist.AddContainerPointsRequest) error {
	return RunOrderWorker(ctx, request.OrderRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *OrderWorkerParams) (err error) {
			return txAddContainerPoints(request, params)
		})
}

func txAddContainerPoints(request dto4logist.AddContainerPointsRequest, params *OrderWorkerParams) error {
	changed := false
	for _, cp := range request.ContainerPoints {
		containerPoint := params.Order.Dto.GetContainerPoint(cp.ContainerID, cp.ShippingPointID)
		if containerPoint == nil {
			params.Order.Dto.ContainerPoints = append(params.Order.Dto.ContainerPoints, &cp)
			changed = true
		} else {
			for _, task := range cp.Tasks {
				if !containerPoint.HasTask(task) {
					containerPoint.Tasks = append(containerPoint.Tasks, task)
					changed = true
				}
			}
		}
	}
	if changed {
		params.Changed.ContainerPoints = true
	}
	return nil
}
