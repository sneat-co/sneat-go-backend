package facade4logist

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// DeleteContainerPoints deletes container point from an order
func DeleteContainerPoints(ctx facade.ContextWithUser, request dto4logist.ContainerPointsRequest) error {
	return RunOrderWorker(ctx, ctx.User(), request.OrderRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) error {
		return txDeleteContainerPoints(ctx, tx, params, request)
	})
}

func txDeleteContainerPoints(_ context.Context, _ dal.ReadwriteTransaction, params *OrderWorkerParams, request dto4logist.ContainerPointsRequest) error {
	orderDto := params.Order.Dto
	_, orderContainer := orderDto.GetContainerByID(request.ContainerID)

	// Check if already deleted
	if orderContainer == nil {
		return nil // Nothing to delete
	}

	if err := deleteContainerPoints(request, orderDto, params); err != nil {
		return err
	}

	if err := deleteContainerSegments(request, orderDto, params); err != nil {
		return err
	}

	return nil
}

func deleteContainerPoints(request dto4logist.ContainerPointsRequest, orderDto *dbo4logist.OrderDbo, params *OrderWorkerParams) error {
	containerPoints := make([]*dbo4logist.ContainerPoint, 0, len(orderDto.ContainerPoints))
	for _, cp := range orderDto.ContainerPoints {
		if cp.ContainerID == request.ContainerID && slice.Index(request.ShippingPointIDs, cp.ShippingPointID) >= 0 {
			continue
		}
		containerPoints = append(containerPoints, cp)
	}
	if len(containerPoints) == len(orderDto.ContainerPoints) {
		return nil // Nothing to delete
	}
	orderDto.ContainerPoints = containerPoints
	params.Changed.ContainerPoints = true
	return nil
}

//func deleteRefsToContainerPointsFromShippingPoints(request dto4logist.ContainerPointsRequest, orderDto *dbo4logist.OrderDbo, params *OrderWorkerParams) error {
//	orderDto.GetShippingPointByID()
//	return nil
//}

func deleteContainerSegments(request dto4logist.ContainerPointsRequest, orderDto *dbo4logist.OrderDbo, params *OrderWorkerParams) error {
	segments := make([]*dbo4logist.ContainerSegment, 0, len(orderDto.Segments))
	for _, segment := range orderDto.Segments {
		if segment.ContainerID == request.ContainerID {
			if slice.Index(request.ShippingPointIDs, segment.From.ShippingPointID) >= 0 ||
				slice.Index(request.ShippingPointIDs, segment.To.ShippingPointID) >= 0 {
				continue
			}
		}
		segments = append(segments, segment)
	}
	orderDto.Segments = segments
	params.Changed.Segments = true
	return nil
}
