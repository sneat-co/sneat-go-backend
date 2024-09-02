package dal4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// ModuleSpaceWorkerParams passes data to a space worker
type ModuleSpaceWorkerParams[D SpaceModuleDbo] struct {
	*SpaceWorkerParams
	SpaceModuleEntry   record.DataWithID[string, D]
	SpaceModuleUpdates []dal.Update
}

func (v *ModuleSpaceWorkerParams[D]) AddSpaceModuleUpdates(updates ...dal.Update) {
	if len(updates) > 0 {
		v.SpaceModuleUpdates = append(v.SpaceModuleUpdates, updates...)
		v.SpaceModuleEntry.Record.MarkAsChanged()
	}
}

func (v *ModuleSpaceWorkerParams[D]) GetRecords(ctx context.Context, tx dal.ReadSession, records ...dal.Record) error {
	return v.SpaceWorkerParams.GetRecords(ctx, tx, append(records, v.SpaceModuleEntry.Record)...)
}

type ModuleDbo interface {
	Validate() error
}

type SpaceModuleDbo = ModuleDbo

func RunModuleSpaceWorkerTx[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if worker == nil {
		panic("worker is nil")
	}
	spaceWorkerParams := NewSpaceWorkerParams(userCtx, request.SpaceID)
	params := NewSpaceModuleWorkerParams(moduleID, spaceWorkerParams, data)
	return runModuleSpaceWorkerReadwriteTx(ctx, tx, params, worker)
}

func NewSpaceModuleWorkerParams[D SpaceModuleDbo](
	moduleID string,
	spaceWorkerParams *SpaceWorkerParams,
	data D,
) *ModuleSpaceWorkerParams[D] {
	return &ModuleSpaceWorkerParams[D]{
		SpaceWorkerParams: spaceWorkerParams,
		SpaceModuleEntry: record.NewDataWithID(moduleID,
			dal.NewKeyWithParentAndID(spaceWorkerParams.Space.Key, dbo4spaceus.SpaceModulesCollection, moduleID),
			data,
		),
	}
}

func runModuleSpaceWorkerReadonlyTx[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleSpaceWorkerParams[D],
	worker func(ctx context.Context, tx dal.ReadTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute module space worker: %w", err)
	}
	return nil
}

func runModuleSpaceWorkerReadwriteTx[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleSpaceWorkerParams[D],
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute module space worker: %w", err)
	}
	if err = applySpaceUpdates(ctx, tx, params.SpaceWorkerParams); err != nil {
		return fmt.Errorf("space module worker failed to apply space record updates: %w", err)
	}
	if err = applySpaceModuleUpdates(ctx, tx, params); err != nil {
		return fmt.Errorf("space module worker failed to apply space module record updates: %w", err)
	}
	return nil
}

func RunReadonlyModuleSpaceWorker[D SpaceModuleDbo](
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	spaceWorkerParams := NewSpaceWorkerParams(userCtx, request.SpaceID)
	params := NewSpaceModuleWorkerParams(moduleID, spaceWorkerParams, data)

	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return runModuleSpaceWorkerReadonlyTx(ctx, tx, params, worker)
	})
}

func RunModuleSpaceWorker[D SpaceModuleDbo](
	ctx context.Context,
	userCtx facade.UserContext,
	spaceID string,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	return RunSpaceWorkerWithUserContext(ctx, userCtx, spaceID,
		func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *SpaceWorkerParams) (err error) {
			params := NewSpaceModuleWorkerParams(moduleID, spaceWorkerParams, data)
			if err = runModuleSpaceWorkerReadwriteTx(ctx, tx, params, worker); err != nil {
				return fmt.Errorf("failed to execute module space worker: %w", err)
			}
			return
		})
}
