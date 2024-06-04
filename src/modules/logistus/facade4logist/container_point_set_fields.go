package facade4logist

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// SetContainerPointFields adds/remove task for a container point
func SetContainerPointFields(ctx context.Context, user facade.User, request dto4logist.SetContainerPointFieldsRequest) error {
	return RunOrderWorker(ctx, user, request.OrderRequest,
		func(_ context.Context, _ dal.ReadwriteTransaction, params *OrderWorkerParams) error {
			return txSetContainerPointFields(params, request)
		},
	)
}

func txSetContainerPointFields(
	params *OrderWorkerParams,
	request dto4logist.SetContainerPointFieldsRequest,
) error {
	containerPoint := params.Order.Dto.GetContainerPoint(request.ContainerID, request.ShippingPointID)
	if containerPoint == nil {
		containerPoint = &dbo4logist.ContainerPoint{
			ContainerID:     request.ContainerID,
			ShippingPointID: request.ShippingPointID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: dbo4logist.ShippingPointStatusPending,
			},
		}
	}
	for name, value := range request.SetStrings {
		switch name {
		case "notes":
			if containerPoint.Notes != value {
				containerPoint.Notes = value
				params.Changed.ContainerPoints = true
			}
		case "refNumber":
			if containerPoint.RefNumber != value {
				containerPoint.RefNumber = value
				params.Changed.ContainerPoints = true
			}
		default:
			return validation.NewErrBadRequestFieldValue("setStrings", "unknown container point field name: "+name)
		}
	}
	return nil
}
