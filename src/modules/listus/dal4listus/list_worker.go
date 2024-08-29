package dal4listus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ListWorkerParams struct {
	dal4spaceus.SpaceWorkerParams
	List        ListEntry
	ListUpdates []dal.Update
}

type ListWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, listWorkerParams *ListWorkerParams) (err error)

func RunListWorker(ctx context.Context, userCtx facade.UserContext, request dto4listus.ListRequest, worker ListWorker) (err error) {
	params := ListWorkerParams{
		SpaceWorkerParams: dal4spaceus.SpaceWorkerParams{
			Space: dbo4spaceus.NewSpaceEntry(request.SpaceID),
		},
		List: NewSpaceListEntry(request.SpaceID, request.ListID),
	}
	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = GetListForUpdate(ctx, tx, params.List); err != nil {
			return err
		}
		if err = worker(ctx, tx, &params); err != nil {
			return err
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
