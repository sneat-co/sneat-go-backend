package facade4logist

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
)

// UpdateContainerPoint updates container point in an order
func UpdateContainerPoint(ctx context.Context, userCtx facade.UserContext, request dto4logist.UpdateContainerPointRequest) error {
	return RunOrderWorker(ctx, userCtx, request.OrderRequest,
		func(_ context.Context, _ dal.ReadwriteTransaction, params *OrderWorkerParams) error {
			return txUpdateContainerPoint(params, request)
		},
	)
}

func txUpdateContainerPoint(
	params *OrderWorkerParams,
	request dto4logist.UpdateContainerPointRequest,
) error {
	orderDto := params.Order.Dto
	containerPoint := orderDto.WithContainerPoints.GetContainerPoint(request.ContainerID, request.ShippingPointID)
	containerPoint.ToLoad = request.ToLoad
	if containerPoint.ToLoad.IsEmpty() {
		containerPoint.ToLoad = nil
	}

	if request.ArrivesDate != nil {
		containerPoint.Arrival.ScheduledDate = *request.ArrivesDate
		segmentKey := dbo4logist.ContainerSegmentKey{
			ContainerID: request.ContainerID,
			To: dbo4logist.SegmentEndpoint{
				ShippingPointID: request.ShippingPointID,
			},
		}
		if segment := orderDto.GetSegmentByKey(segmentKey); segment != nil {
			params.Changed.Segments = true
			if segment.Dates == nil {
				segment.Dates = &dbo4logist.SegmentDates{}
			}
			segment.Dates.Arrives = containerPoint.Arrival.ScheduledDate
		}
	}
	if request.DepartsDate != nil {
		containerPoint.Departure.ScheduledDate = *request.DepartsDate
		segmentKey := dbo4logist.ContainerSegmentKey{
			ContainerID: request.ContainerID,
			From: dbo4logist.SegmentEndpoint{
				ShippingPointID: request.ShippingPointID,
			},
		}
		if segment := orderDto.GetSegmentByKey(segmentKey); segment != nil {
			params.Changed.Segments = true
			if segment.Dates == nil {
				segment.Dates = &dbo4logist.SegmentDates{}
			}
			segment.Dates.Departs = containerPoint.Departure.ScheduledDate
		}
	}
	params.Changed.ContainerPoints = true
	return nil
}
