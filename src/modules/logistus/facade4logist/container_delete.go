package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
)

// DeleteContainer deletes container from an order
func DeleteContainer(ctx context.Context, userCtx facade.UserContext, request dto4logist.ContainerRequest) error {
	err := RunOrderWorker(ctx, userCtx, request.OrderRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) error {
		return deleteContainerTx(request, params)
	})
	if err != nil {
		return fmt.Errorf("failed to delete container with id=[%s]: %w", request.ContainerID, err)
	}
	return nil
}

func deleteContainerTx(request dto4logist.ContainerRequest, params *OrderWorkerParams) error {
	orderDto := params.Order.Dto
	var containerFound bool
	orderDto.Containers, containerFound = orderDto.RemoveContainer(request.ContainerID)
	if !containerFound {
		return nil
	}
	segmentsToDelete := orderDto.GetSegmentsByFilter(dbo4logist.SegmentsFilter{ContainerIDs: []string{request.ContainerID}})
	orderDto.Segments = orderDto.DeleteSegments(segmentsToDelete)
	orderDto.ContainerPoints = orderDto.RemoveContainerPointsByContainerID(request.ContainerID)
	params.Changed.Containers = true
	params.Changed.ContainerPoints = true
	return nil
}
