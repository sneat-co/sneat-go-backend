package dal4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

type SpaceItemDbo = interface {
	Validate() error
}

// SpaceItemRequest DTO
type SpaceItemRequest struct {
	dto4spaceus.SpaceRequest
	ID string `json:"id"`
}

// Validate returns error if not valid
func (v SpaceItemRequest) Validate() error {
	if v.ID == "" {
		return validation.NewErrRequestIsMissingRequiredField("id")
	}
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	return nil
}

// SliceIndexes DTO
type SliceIndexes struct {
	Start int
	End   int
}

type Brief = interface {
	Validate()
}

type BriefsAdapter[ModuleDbo SpaceModuleDbo] interface {
	DeleteBrief(space ModuleDbo, id string) ([]dal.Update, error)
	GetBriefsCount(team ModuleDbo) int
}

type mapBriefsAdapter[ModuleDbo SpaceModuleDbo] struct {
	getBriefsCount func(space ModuleDbo) int
	deleteBrief    func(space ModuleDbo, id string) ([]dal.Update, error)
}

func (v mapBriefsAdapter[ModuleDbo]) DeleteBrief(teamModuleDbo ModuleDbo, id string) ([]dal.Update, error) {
	return v.deleteBrief(teamModuleDbo, id)
}

func (v mapBriefsAdapter[ModuleDbo]) GetBriefsCount(teamModuleDbo ModuleDbo) int {
	return v.getBriefsCount(teamModuleDbo)
}

func NewMapBriefsAdapter[ModuleDbo SpaceModuleDbo](
	getBriefsCount func(teamModuleDbo ModuleDbo) int,
	deleteBrief func(teamModuleDbo ModuleDbo, id string) ([]dal.Update, error),
) BriefsAdapter[ModuleDbo] {
	return mapBriefsAdapter[ModuleDbo]{
		getBriefsCount: getBriefsCount,
		deleteBrief:    deleteBrief,
	}
}

// SpaceItemWorkerParams defines params for team item worker
type SpaceItemWorkerParams[ModuleDbo SpaceModuleDbo, ItemDbo SpaceItemDbo] struct {
	*ModuleSpaceWorkerParams[ModuleDbo]
	SpaceItem        record.DataWithID[string, ItemDbo]
	SpaceItemUpdates []dal.Update
}

func RunSpaceItemWorker[ModuleDbo SpaceModuleDbo, ItemDbo SpaceItemDbo](
	ctx context.Context,
	userCtx facade.UserContext,
	request SpaceItemRequest,
	moduleID string,
	spaceModuleData ModuleDbo,
	spaceItemCollection string,
	spaceItemDbo ItemDbo,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) (err error),
) (err error) {
	return RunModuleSpaceWorker(ctx, userCtx, request.SpaceID, moduleID, spaceModuleData,
		func(ctx context.Context, tx dal.ReadwriteTransaction, moduleSpaceWorkerParams *ModuleSpaceWorkerParams[ModuleDbo]) (err error) {
			teamItemKey := dal.NewKeyWithParentAndID(moduleSpaceWorkerParams.SpaceModuleEntry.Key, spaceItemCollection, request.ID)
			params := SpaceItemWorkerParams[ModuleDbo, ItemDbo]{
				ModuleSpaceWorkerParams: moduleSpaceWorkerParams,
				SpaceItem:               record.NewDataWithID(request.ID, teamItemKey, spaceItemDbo),
			}
			if err = worker(ctx, tx, &params); err != nil {
				return err
			}
			if len(params.SpaceItemUpdates) > 0 {
				if err = tx.Update(ctx, teamItemKey, params.SpaceItemUpdates); err != nil {
					return fmt.Errorf("failed to update team item record: %w", err)
				}
			}
			return nil
		},
	)
}

// DeleteSpaceItem deletes team item
func DeleteSpaceItem[ModuleDbo SpaceModuleDbo, ItemDbo SpaceItemDbo](
	ctx context.Context,
	userCtx facade.UserContext,
	request SpaceItemRequest,
	moduleID string,
	moduleData ModuleDbo,
	teamItemCollection string,
	teamItemDbo ItemDbo,
	briefsAdapter BriefsAdapter[ModuleDbo],
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) (err error),
) (err error) {
	return RunSpaceItemWorker(ctx, userCtx, request, moduleID, moduleData, teamItemCollection, teamItemDbo,
		func(ctx context.Context, tx dal.ReadwriteTransaction, teamItemWorkerParams *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) (err error) {
			return deleteSpaceItemTxWorker[ModuleDbo](ctx, tx, teamItemWorkerParams, briefsAdapter, worker)
		},
	)
}

func deleteSpaceItemTxWorker[ModuleDbo SpaceModuleDbo, ItemDbo SpaceItemDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *SpaceItemWorkerParams[ModuleDbo, ItemDbo],
	briefsAdapter BriefsAdapter[ModuleDbo],
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, p *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) error,
) (err error) {
	if err = tx.Get(ctx, params.Space.Record); err != nil {
		return
	}
	if err = tx.Get(ctx, params.SpaceItem.Record); err != nil && !dal.IsNotFound(err) {
		return err
	}
	if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
		return
	}
	if worker != nil {
		if err = worker(ctx, tx, params); err != nil {
			return fmt.Errorf("failed to execute teamItemWorker: %w", err)
		}
	}
	//if err = decrementCounter(&params); err != nil {
	//	return err
	//}
	if len(params.SpaceUpdates) > 0 {
		if err = TxUpdateSpace(ctx, tx, params.Started, params.Space, params.SpaceUpdates); err != nil {
			return fmt.Errorf("failed to update team record: %w", err)
		}
	}
	var spaceModuleUpdates []dal.Update
	if spaceModuleUpdates, err = deleteBrief[ModuleDbo](params.SpaceModuleEntry, params.SpaceItem.ID, briefsAdapter, params.SpaceModuleUpdates); err != nil {
		return err
	} else {
		params.AddSpaceModuleUpdates(spaceModuleUpdates...)
	}

	if params.SpaceItem.Record.Exists() {
		if err = tx.Delete(ctx, params.SpaceItem.Key); err != nil {
			return fmt.Errorf("failed to delete team item record by key=%v: %w", params.SpaceItem.Key, err)
		}
	}
	return err
}

func deleteBrief[D SpaceModuleDbo](teamModuleEntry record.DataWithID[string, D], itemID string, adapter BriefsAdapter[D], updates []dal.Update) ([]dal.Update, error) {
	if adapter == nil {
		return updates, nil
	}
	return adapter.DeleteBrief(teamModuleEntry.Data, itemID)
}
