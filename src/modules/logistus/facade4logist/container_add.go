package facade4logist

import (
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// AddContainers adds containers to an order
func AddContainers(ctx facade.ContextWithUser, request dto4logist.AddContainersRequest) error {
	return RunOrderWorker(ctx, request.OrderRequest,
		func(_ facade.ContextWithUser, _ dal.ReadwriteTransaction, params *OrderWorkerParams) error {
			return addContainersTx(params, request)
		})
}

func addContainersTx(params *OrderWorkerParams, request dto4logist.AddContainersRequest) error {
	if params.Order.Dto == nil {
		panic("params.Order.Data == nil")
	}
	for _, c := range request.Containers {
		containerID := params.Order.Dto.GenerateRandomContainerID()
		containerBrief := dbo4logist.OrderContainer{
			ID:                 containerID,
			OrderContainerBase: c.OrderContainerBase,
		}
		params.Order.Dto.Containers = append(params.Order.Dto.Containers, &containerBrief)
		for j, point := range c.Points {
			_, shippingPoint := params.Order.Dto.GetShippingPointByID(point.ShippingPointID)
			if shippingPoint == nil {
				return validation.NewErrBadRequestFieldValue(fmt.Sprintf("shippingPointIDs[%v]", j), "unknown shipping point ContactID")
			}
			containerPoint := dbo4logist.ContainerPoint{
				ContainerID:       containerID,
				ShippingPointID:   point.ShippingPointID,
				ShippingPointBase: shippingPoint.ShippingPointBase,
			}
			params.Order.Dto.ContainerPoints = append(params.Order.Dto.ContainerPoints, &containerPoint)
		}
	}
	params.Changed.Containers = true
	params.Changed.ContainerPoints = true
	return nil
}
