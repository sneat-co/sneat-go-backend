package facade4logist

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// UpdateShippingPoint updates shipping point in an order
func UpdateShippingPoint(ctx facade.ContextWithUser, request dto4logist.UpdateShippingPointRequest) error {
	return RunOrderWorker(ctx, request.OrderRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *OrderWorkerParams) error {
			return txUpdateShippingPoint(ctx, tx, params, request)
		})
}

func txUpdateShippingPoint(_ context.Context, _ dal.ReadwriteTransaction, params *OrderWorkerParams, request dto4logist.UpdateShippingPointRequest) error {
	orderDto := params.Order.Dto

	_, shippingPoint := orderDto.GetShippingPointByID(request.ShippingPointID)

	for name, value := range request.SetDates {
		switch name {
		case "scheduledStartDate":
			shippingPoint.ScheduledStartDate = value
		case "scheduledEndDate":
			shippingPoint.ScheduledEndDate = value
		default:
			return validation.NewErrBadRequestFieldValue("setDates."+name, "unknown field name: "+name)
		}

	}
	for name, value := range request.SetStrings {
		switch name {
		case "notes":
			shippingPoint.Notes = value
		default:
			return validation.NewErrBadRequestFieldValue("setStrings."+name, "unknown field name: "+name)
		}
	}
	params.Changed.ShippingPoints = true
	return nil
}
