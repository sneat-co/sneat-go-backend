package facade4logist

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-core/facade"
)

// SetContainerPointFreightFields adds/remove task for a container point
func SetContainerPointFreightFields(ctx context.Context, user facade.User, request dto4logist.SetContainerPointFreightFieldsRequest) error {
	return RunOrderWorker(ctx, user, request.OrderRequest,
		func(_ context.Context, _ dal.ReadwriteTransaction, params *OrderWorkerParams) error {
			return txSetContainerPointFreightFields(params, request)
		},
	)
}

func txSetContainerPointFreightFields(
	params *OrderWorkerParams,
	request dto4logist.SetContainerPointFreightFieldsRequest,
) error {
	containerPoint := params.Order.Dto.GetContainerPoint(request.ContainerID, request.ShippingPointID)
	if containerPoint == nil {
		return errors.New("container point not found byt containerID & shippingPointID")
	}
	setNumber := func(freightLoad *models4logist.FreightLoad) (*models4logist.FreightLoad, error) {
		if freightLoad == nil {
			freightLoad = &models4logist.FreightLoad{}
		}
		for name, value := range request.Integers {
			switch name {
			case "numberOfPallets":
				if value == nil {
					freightLoad.NumberOfPallets = 0
				} else {
					freightLoad.NumberOfPallets = *value
				}
			case "grossWeightKg":
				if value == nil {
					freightLoad.GrossWeightKg = 0
				} else {
					freightLoad.GrossWeightKg = *value
				}
			case "volumeM3":
				if value == nil {
					freightLoad.VolumeM3 = 0
				} else {
					freightLoad.VolumeM3 = *value
				}
			default:
				return freightLoad, errors.New("unknown freight load name: " + name)
			}
		}
		return freightLoad, nil
	}
	var err error = nil
	if len(request.Integers) > 0 {
		switch request.Task {
		case models4logist.ShippingPointTaskLoad:
			containerPoint.ToLoad, err = setNumber(containerPoint.ToLoad)
		case models4logist.ShippingPointTaskUnload:
			containerPoint.ToUnload, err = setNumber(containerPoint.ToUnload)
		}
	}
	if err != nil {
		return err
	}
	params.Changed.ContainerPoints = true
	return nil
}
