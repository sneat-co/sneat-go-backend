package facade4logist

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
)

// SetContainerFields sets container number in an order
func SetContainerFields(ctx context.Context, user facade.User, request dto4logist.SetContainerFieldsRequest) error {
	return RunOrderWorker(ctx, user, request.OrderRequest,
		func(_ context.Context, _ dal.ReadwriteTransaction, params *OrderWorkerParams) error {
			return txSetContainerFields(params, request)
		},
	)
}

func txSetContainerFields(
	params *OrderWorkerParams,
	request dto4logist.SetContainerFieldsRequest,
) error {
	_, container := params.Order.Dto.GetContainerByID(request.ContainerID)
	if container == nil {
		return validation.NewErrBadRequestFieldValue("containerID", "container not found")
	}
	for name, value := range request.SetStrings {
		switch name {
		case "number":
			container.Number = strings.TrimSpace(value)
		case "instructions":
			container.Instructions = strings.TrimSpace(value)
		default:
			return validation.NewErrBadRequestFieldValue("setStrings", "unknown container point field name: "+name)
		}
	}
	params.Changed.Containers = true
	return nil
}
