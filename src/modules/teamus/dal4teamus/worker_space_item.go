package dal4teamus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

type SpaceItemDbo = interface {
	Validate() error
}

// SpaceItemRequest DTO
type SpaceItemRequest struct {
	dto4teamus.SpaceRequest
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
	user facade.User,
	request SpaceItemRequest,
	moduleID string,
	spaceModuleData ModuleDbo,
	spaceItemCollection string,
	spaceItemDbo ItemDbo,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) (err error),
) (err error) {
	return RunModuleSpaceWorker(ctx, user, request.SpaceRequest, moduleID, spaceModuleData,
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

//var RunSpaceItemWorker = func(ctx context.Context, db dal.DB, request TeamItemRunnerInput, worker teamItemWorker) (err error) {
//	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
//		params := SpaceItemWorkerParams{
//			Started: time.Now(),
//		}
//		if request.IsTeamRecordNeeded {
//			params.TeamKey = dal.NewKeyWithID(dbo4teamus.SpacesCollection, request.InviteID)
//			params.SpaceID = new(dbo4teamus.TeamDto)
//			teamRecord := dal.NewRecordWithData(params.TeamKey, params.SpaceID)
//			if err = tx.Get(ctx, teamRecord); err != nil {
//				return
//			}
//		}
//		if err = tx.Get(ctx, request.SpaceItem); err != nil {
//			return err
//		}
//		err = worker(c, tx, &params)
//		if err != nil {
//			return fmt.Errorf("failed to execute spaceWorker: %w", err)
//		}
//		if len(params.SpaceUpdates) > 0 {
//			if request.IsTeamRecordNeeded {
//				if err = TxUpdateSpace(ctx, tx, params.Started, params.SpaceID, params.TeamKey, params.SpaceUpdates); err != nil {
//					return fmt.Errorf("failed to update team record: %w", err)
//				}
//			} else {
//				return fmt.Errorf("got %v team updates but request indicated team record is not required", len(params.SpaceUpdates))
//			}
//		}
//		return err
//	})
//}

//func incrementCounter(params *SpaceWorkerParams, moduleID, counter string) (err error) {
//	if counter == "" {
//		return nil
//	}
//	numberFieldName := "numberOf." + counter
//	for _, teamUpdate := range params.SpaceUpdates {
//		if teamUpdate.Field == numberFieldName {
//			return
//		}
//	}
//	numberOfItems, hasNumber := params.Space.Data.NumberOf[counter]
//	if hasNumber {
//		if numberOfItems >= 0 {
//			params.SpaceUpdates = append(params.SpaceUpdates, dal.Update{Field: numberFieldName, Value: dal.Increment(1)})
//		} else {
//			return validation.NewErrBadRecordFieldValue(numberFieldName, fmt.Sprintf("has negative value: %d", numberOfItems))
//		}
//	} else {
//		params.SpaceUpdates = append(params.SpaceUpdates, dal.Update{Field: numberFieldName, Value: 1})
//	}
//	return nil
//}
//
//func decrementCounter(params *SpaceItemWorkerParams) (err error) {
//	if params.Counter != "" {
//		numberFieldName := "numberOf." + params.Counter
//		for _, teamUpdate := range params.SpaceUpdates {
//			if teamUpdate.Field == numberFieldName {
//				return
//			}
//		}
//		numberOfItems, hasNumber := params.Space.Data.NumberOf[params.Counter]
//		if hasNumber {
//			if numberOfItems > 0 {
//				params.SpaceUpdates = append(params.SpaceUpdates, dal.Update{Field: numberFieldName, Value: dal.Increment(-1)})
//			} else {
//				return validation.NewErrBadRecordFieldValue(numberFieldName, fmt.Sprintf("an attempt to decrement counter with negative or zero value: %d", numberOfItems))
//			}
//		} else {
//			return validation.NewErrBadRecordFieldValue(numberFieldName, "an attempt to decrement non existing counter")
//		}
//	}
//	return nil
//}

// DeleteSpaceItem deletes team item
func DeleteSpaceItem[ModuleDbo SpaceModuleDbo, ItemDbo SpaceItemDbo](
	ctx context.Context,
	user facade.User,
	request SpaceItemRequest,
	moduleID string,
	moduleData ModuleDbo,
	teamItemCollection string,
	teamItemDbo ItemDbo,
	briefsAdapter BriefsAdapter[ModuleDbo],
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) (err error),
) (err error) {
	return RunSpaceItemWorker(ctx, user, request, moduleID, moduleData, teamItemCollection, teamItemDbo,
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
	var teamModuleUpdates []dal.Update
	if teamModuleUpdates, err = deleteBrief[ModuleDbo](params.SpaceModuleEntry, params.SpaceItem.ID, briefsAdapter, params.SpaceModuleUpdates); err != nil {
		return err
	} else {
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, teamModuleUpdates...)

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
