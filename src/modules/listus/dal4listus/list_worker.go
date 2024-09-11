package dal4listus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ListWorkerParams struct {
	*dal4spaceus.ModuleSpaceWorkerParams[*dbo4listus.ListusSpaceDbo]
	List        ListEntry
	ListUpdates []dal.Update
}

type ListWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, listWorkerParams *ListWorkerParams) (err error)

func RunListWorker(ctx context.Context, userCtx facade.UserContext, request dto4listus.ListRequest, worker ListWorker) (err error) {
	err = dal4spaceus.RunModuleSpaceWorker(ctx, userCtx, request.SpaceID, "listus", new(dbo4listus.ListusSpaceDbo), func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *dal4spaceus.ModuleSpaceWorkerParams[*dbo4listus.ListusSpaceDbo]) (err error) {
		params := ListWorkerParams{
			ModuleSpaceWorkerParams: spaceWorkerParams,
			List:                    NewListEntry(request.SpaceID, request.ListID),
		}
		if err = GetListForUpdate(ctx, tx, params.List); err != nil {
			if dal.IsNotFound(err) && dbo4listus.IsStandardList(request.ListID) {
				// It's OK to miss a standard list record - should be created automatically
			} else {
				return err
			}
		}
		if err = worker(ctx, tx, &params); err != nil {
			return err
		}
		if params.List.Data.Title == params.List.ID && params.List.Record.Exists() {
			params.List.Data.Title = ""
			params.ListUpdates = append(params.ListUpdates, dal.Update{
				Field: "title",
				Value: "",
			})
			params.List.Record.MarkAsChanged()
		}
		if updateCount := len(params.ListUpdates); updateCount > 0 {
			if !params.List.Record.HasChanged() {
				err = fmt.Errorf("got %d list updates but list record is not marked as changed", updateCount)
				return
			}
			if err = tx.Update(ctx, params.List.Record.Key(), params.ListUpdates); err != nil {
				return
			}
		}

		return
	})
	return
}
