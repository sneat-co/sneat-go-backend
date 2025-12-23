package facade4logist

import (
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteSegments deletes segments from an order
func DeleteSegments(ctx facade.ContextWithUser, request dto4logist.DeleteSegmentsRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return RunOrderWorker(ctx, request.OrderRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *OrderWorkerParams) error {
			return deleteSegments(params, request)
		})
}

func deleteSegments(params *OrderWorkerParams, request dto4logist.DeleteSegmentsRequest) error {
	if err := request.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	orderDto := params.Order.Dto

	segmentsToDelete := orderDto.GetSegmentsByFilter(request.SegmentsFilter)
	if len(segmentsToDelete) > 0 {
		orderDto.Segments = orderDto.DeleteSegments(segmentsToDelete)
		params.Changed.Segments = true

		for _, segment := range segmentsToDelete {
			fromContainerPoint := orderDto.GetContainerPoint(segment.ContainerID, segment.From.ShippingPointID)
			if fromContainerPoint != nil && fromContainerPoint.Departure != nil {
				fromContainerPoint.Departure.ScheduledDate = ""
				params.Changed.ContainerPoints = true
			}
			toContainerPoint := orderDto.GetContainerPoint(segment.ContainerID, segment.To.ShippingPointID)
			if toContainerPoint != nil && toContainerPoint.Arrival != nil {
				toContainerPoint.Arrival.ScheduledDate = ""
				params.Changed.ContainerPoints = true
			}
		}
	}

	return nil
}
