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

type TeamItemDbo = interface {
	Validate() error
}

// TeamItemRequest DTO
type TeamItemRequest struct {
	dto4teamus.TeamRequest
	ID string `json:"id"`
}

// Validate returns error if not valid
func (v TeamItemRequest) Validate() error {
	if v.ID == "" {
		return validation.NewErrRequestIsMissingRequiredField("id")
	}
	if err := v.TeamRequest.Validate(); err != nil {
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

type BriefsAdapter[ModuleDbo TeamModuleDbo] interface {
	DeleteBrief(team ModuleDbo, id string) ([]dal.Update, error)
	GetBriefsCount(team ModuleDbo) int
}

type mapBriefsAdapter[ModuleDbo TeamModuleDbo] struct {
	getBriefsCount func(team ModuleDbo) int
	deleteBrief    func(team ModuleDbo, id string) ([]dal.Update, error)
}

func (v mapBriefsAdapter[ModuleDbo]) DeleteBrief(teamModuleDbo ModuleDbo, id string) ([]dal.Update, error) {
	return v.deleteBrief(teamModuleDbo, id)
}

func (v mapBriefsAdapter[ModuleDbo]) GetBriefsCount(teamModuleDbo ModuleDbo) int {
	return v.getBriefsCount(teamModuleDbo)
}

func NewMapBriefsAdapter[ModuleDbo TeamModuleDbo](
	getBriefsCount func(teamModuleDbo ModuleDbo) int,
	deleteBrief func(teamModuleDbo ModuleDbo, id string) ([]dal.Update, error),
) BriefsAdapter[ModuleDbo] {
	return mapBriefsAdapter[ModuleDbo]{
		getBriefsCount: getBriefsCount,
		deleteBrief:    deleteBrief,
	}
}

// TeamItemWorkerParams defines params for team item worker
type TeamItemWorkerParams[ModuleDbo TeamModuleDbo, ItemDbo TeamItemDbo] struct {
	*ModuleTeamWorkerParams[ModuleDbo]
	TeamItem        record.DataWithID[string, ItemDbo]
	TeamItemUpdates []dal.Update
}

func RunTeamItemWorker[ModuleDbo TeamModuleDbo, ItemDbo TeamItemDbo](
	ctx context.Context,
	user facade.User,
	request TeamItemRequest,
	moduleID string,
	teamModuleData ModuleDbo,
	teamItemCollection string,
	teamItemDbo ItemDbo,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *TeamItemWorkerParams[ModuleDbo, ItemDbo]) (err error),
) (err error) {
	return RunModuleTeamWorker(ctx, user, request.TeamRequest, moduleID, teamModuleData,
		func(ctx context.Context, tx dal.ReadwriteTransaction, moduleTeamWorkerParams *ModuleTeamWorkerParams[ModuleDbo]) (err error) {
			teamItemKey := dal.NewKeyWithParentAndID(moduleTeamWorkerParams.TeamModuleEntry.Key, teamItemCollection, request.ID)
			params := TeamItemWorkerParams[ModuleDbo, ItemDbo]{
				ModuleTeamWorkerParams: moduleTeamWorkerParams,
				TeamItem:               record.NewDataWithID(request.ID, teamItemKey, teamItemDbo),
			}
			if err = worker(ctx, tx, &params); err != nil {
				return err
			}
			if len(params.TeamItemUpdates) > 0 {
				if err = tx.Update(ctx, teamItemKey, params.TeamItemUpdates); err != nil {
					return fmt.Errorf("failed to update team item record: %w", err)
				}
			}
			return nil
		},
	)
}

//var RunTeamItemWorker = func(ctx context.Context, db dal.DB, request TeamItemRunnerInput, worker teamItemWorker) (err error) {
//	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
//		params := TeamItemWorkerParams{
//			Started: time.Now(),
//		}
//		if request.IsTeamRecordNeeded {
//			params.TeamKey = dal.NewKeyWithID(dbo4teamus.TeamsCollection, request.InviteID)
//			params.TeamID = new(dbo4teamus.TeamDto)
//			teamRecord := dal.NewRecordWithData(params.TeamKey, params.TeamID)
//			if err = tx.Get(ctx, teamRecord); err != nil {
//				return
//			}
//		}
//		if err = tx.Get(ctx, request.TeamItem); err != nil {
//			return err
//		}
//		err = worker(c, tx, &params)
//		if err != nil {
//			return fmt.Errorf("failed to execute teamWorker: %w", err)
//		}
//		if len(params.TeamUpdates) > 0 {
//			if request.IsTeamRecordNeeded {
//				if err = TxUpdateTeam(ctx, tx, params.Started, params.TeamID, params.TeamKey, params.TeamUpdates); err != nil {
//					return fmt.Errorf("failed to update team record: %w", err)
//				}
//			} else {
//				return fmt.Errorf("got %v team updates but request indicated team record is not required", len(params.TeamUpdates))
//			}
//		}
//		return err
//	})
//}

//func incrementCounter(params *TeamWorkerParams, moduleID, counter string) (err error) {
//	if counter == "" {
//		return nil
//	}
//	numberFieldName := "numberOf." + counter
//	for _, teamUpdate := range params.TeamUpdates {
//		if teamUpdate.Field == numberFieldName {
//			return
//		}
//	}
//	numberOfItems, hasNumber := params.Team.Data.NumberOf[counter]
//	if hasNumber {
//		if numberOfItems >= 0 {
//			params.TeamUpdates = append(params.TeamUpdates, dal.Update{Field: numberFieldName, Value: dal.Increment(1)})
//		} else {
//			return validation.NewErrBadRecordFieldValue(numberFieldName, fmt.Sprintf("has negative value: %d", numberOfItems))
//		}
//	} else {
//		params.TeamUpdates = append(params.TeamUpdates, dal.Update{Field: numberFieldName, Value: 1})
//	}
//	return nil
//}
//
//func decrementCounter(params *TeamItemWorkerParams) (err error) {
//	if params.Counter != "" {
//		numberFieldName := "numberOf." + params.Counter
//		for _, teamUpdate := range params.TeamUpdates {
//			if teamUpdate.Field == numberFieldName {
//				return
//			}
//		}
//		numberOfItems, hasNumber := params.Team.Data.NumberOf[params.Counter]
//		if hasNumber {
//			if numberOfItems > 0 {
//				params.TeamUpdates = append(params.TeamUpdates, dal.Update{Field: numberFieldName, Value: dal.Increment(-1)})
//			} else {
//				return validation.NewErrBadRecordFieldValue(numberFieldName, fmt.Sprintf("an attempt to decrement counter with negative or zero value: %d", numberOfItems))
//			}
//		} else {
//			return validation.NewErrBadRecordFieldValue(numberFieldName, "an attempt to decrement non existing counter")
//		}
//	}
//	return nil
//}

// DeleteTeamItem deletes team item
func DeleteTeamItem[ModuleDbo TeamModuleDbo, ItemDbo TeamItemDbo](
	ctx context.Context,
	user facade.User,
	request TeamItemRequest,
	moduleID string,
	moduleData ModuleDbo,
	teamItemCollection string,
	teamItemDbo ItemDbo,
	briefsAdapter BriefsAdapter[ModuleDbo],
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *TeamItemWorkerParams[ModuleDbo, ItemDbo]) (err error),
) (err error) {
	return RunTeamItemWorker(ctx, user, request, moduleID, moduleData, teamItemCollection, teamItemDbo,
		func(ctx context.Context, tx dal.ReadwriteTransaction, teamItemWorkerParams *TeamItemWorkerParams[ModuleDbo, ItemDbo]) (err error) {
			return deleteTeamItemTxWorker[ModuleDbo](ctx, tx, teamItemWorkerParams, briefsAdapter, worker)
		},
	)
}

func deleteTeamItemTxWorker[ModuleDbo TeamModuleDbo, ItemDbo TeamItemDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *TeamItemWorkerParams[ModuleDbo, ItemDbo],
	briefsAdapter BriefsAdapter[ModuleDbo],
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, p *TeamItemWorkerParams[ModuleDbo, ItemDbo]) error,
) (err error) {
	if err = tx.Get(ctx, params.Team.Record); err != nil {
		return
	}
	if err = tx.Get(ctx, params.TeamItem.Record); err != nil && !dal.IsNotFound(err) {
		return err
	}
	if err = tx.Get(ctx, params.TeamModuleEntry.Record); err != nil {
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
	if len(params.TeamUpdates) > 0 {
		if err = TxUpdateTeam(ctx, tx, params.Started, params.Team, params.TeamUpdates); err != nil {
			return fmt.Errorf("failed to update team record: %w", err)
		}
	}
	var teamModuleUpdates []dal.Update
	if teamModuleUpdates, err = deleteBrief[ModuleDbo](params.TeamModuleEntry, params.TeamItem.ID, briefsAdapter, params.TeamModuleUpdates); err != nil {
		return err
	} else {
		params.TeamModuleUpdates = append(params.TeamModuleUpdates, teamModuleUpdates...)

	}

	if params.TeamItem.Record.Exists() {
		if err = tx.Delete(ctx, params.TeamItem.Key); err != nil {
			return fmt.Errorf("failed to delete team item record by key=%v: %w", params.TeamItem.Key, err)
		}
	}
	return err
}

func deleteBrief[D TeamModuleDbo](teamModuleEntry record.DataWithID[string, D], itemID string, adapter BriefsAdapter[D], updates []dal.Update) ([]dal.Update, error) {
	if adapter == nil {
		return updates, nil
	}
	return adapter.DeleteBrief(teamModuleEntry.Data, itemID)
}
