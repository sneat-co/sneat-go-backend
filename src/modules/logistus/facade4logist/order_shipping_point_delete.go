package facade4logist

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteShippingPoint deletes shipping point from an order
func DeleteShippingPoint(ctx context.Context, user facade.User, request dto4logist.OrderShippingPointRequest) error {
	return RunOrderWorker(ctx, user, request.OrderRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) error {
		return txDeleteShippingPoint(ctx, tx, params, request)
	})
}

func txDeleteShippingPoint(_ context.Context, _ dal.ReadwriteTransaction, params *OrderWorkerParams, request dto4logist.OrderShippingPointRequest) error {
	orderDto := params.Order.Dto

	var contactID string
	//var counterpartyRole models4logist.Role

	{ // Remove shipping point from order
		shippingPoints := make([]*models4logist.OrderShippingPoint, 0, len(orderDto.ShippingPoints))
		for _, sp := range orderDto.ShippingPoints {
			if sp.ID == request.ShippingPointID {
				if sp.Location != nil {
					contactID = sp.Location.ContactID
				}
				continue
			}
			shippingPoints = append(shippingPoints, sp)
		}
		if len(shippingPoints) == len(orderDto.ShippingPoints) {
			return nil // Nothing to delete
		}
		orderDto.ShippingPoints = shippingPoints
		params.Changed.ShippingPoints = true
	}

	if contactID != "" {
		deleteCounterpartyAndChildren(params, models4logist.CounterpartyRoleDispatchPoint, contactID) // TODO: Why role is hardcoded?
	}

	{ // Remove segments related to the deleted shipping point
		segments := make([]*models4logist.ContainerSegment, 0, len(orderDto.Segments))
		for _, segment := range orderDto.Segments {
			if segment.ContainerSegmentKey.From.ShippingPointID == request.ShippingPointID || segment.ContainerSegmentKey.To.ShippingPointID == request.ShippingPointID {
				continue
			}
			segments = append(segments, segment)
		}
		if len(segments) != len(orderDto.Segments) {
			orderDto.Segments = segments
			params.Changed.Segments = true
		}
	}

	{ // Remove container points related to the deleted shipping point
		containerPoints := make([]*models4logist.ContainerPoint, 0, len(orderDto.ContainerPoints))
		for _, cp := range orderDto.ContainerPoints {
			if cp.ShippingPointID == request.ShippingPointID {
				continue
			}
			containerPoints = append(containerPoints, cp)
		}
		if len(containerPoints) != len(orderDto.ContainerPoints) {
			orderDto.ContainerPoints = containerPoints
			params.Changed.ContainerPoints = true
		}
	}

	return nil
}
