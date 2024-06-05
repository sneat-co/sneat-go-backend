package dal4teamus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
	"time"
)

type teamItemWorker = func(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *TeamItemWorkerParams,
) (err error)

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

// BriefsAdapter defines brief adapters
type BriefsAdapter[D TeamModuleData] struct {
	BriefsFieldName string
	BriefsValue     func(team D) interface{}
	GetBriefsCount  func(team D) int
	GetBriefItemID  func(team D, i int) (id string)
	ShiftBriefs     func(team D, from SliceIndexes, end SliceIndexes)
	TrimBriefs      func(team D, count int)
}

// TeamItemRunnerInput request
type TeamItemRunnerInput[D TeamModuleData] struct {
	dto4teamus.TeamRequest
	IsTeamRecordNeeded bool
	TeamItem           dal.Record
	ShortID            string
	Counter            string
	*BriefsAdapter[D]
}

// Validate validates request
func (v TeamItemRunnerInput[D]) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.ShortID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("shortID")
	}
	return nil
}

// TeamItemWorkerParams defines params for team item worker
type TeamItemWorkerParams struct {
	//Counter     string
	Started     time.Time
	Team        TeamContext
	TeamKey     *dal.Key
	TeamUpdates []dal.Update
	TeamItem    dal.Record
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
func DeleteTeamItem[D TeamModuleData](
	ctx context.Context,
	user facade.User,
	input TeamItemRunnerInput[D],
	moduleID string,
	data D,
	worker teamItemWorker,
) (err error) {
	if err := input.Validate(); err != nil {
		return err
	}
	if input.Counter != "" && worker == nil {
		return validation.NewErrBadRequestFieldValue("counter", "input specifies counter but worker was not provided")
	}
	return RunModuleTeamWorker(ctx, user, input.TeamRequest, moduleID, data,
		func(ctx context.Context, tx dal.ReadwriteTransaction, moduleWorkerParams *ModuleTeamWorkerParams[D]) (err error) {
			params := TeamItemWorkerParams{
				Started:  time.Now(),
				TeamItem: input.TeamItem,
				//Counter:  "",
			}
			if input.IsTeamRecordNeeded {
				params.Team = NewTeamContext(input.TeamID)
				if err = tx.Get(ctx, params.Team.Record); err != nil {
					return
				}
			}
			if worker != nil {
				if err = tx.Get(ctx, input.TeamItem); err != nil {
					return err
				}
				err = worker(ctx, tx, &params)
				if err != nil {
					return fmt.Errorf("failed to execute teamItemWorker: %w", err)
				}
			}
			//if err = decrementCounter(&params); err != nil {
			//	return err
			//}
			if len(moduleWorkerParams.TeamUpdates) > 0 {
				if input.IsTeamRecordNeeded {
					if err = TxUpdateTeam(ctx, tx, moduleWorkerParams.Started, moduleWorkerParams.Team, moduleWorkerParams.TeamUpdates); err != nil {
						return fmt.Errorf("failed to update team record: %w", err)
					}
				} else {
					return fmt.Errorf("got %d team updates but input indicated team record is not required", len(params.TeamUpdates))
				}
			}
			moduleWorkerParams.TeamModuleUpdates = deleteBrief[D](moduleWorkerParams.TeamModuleEntry, input.ShortID, input.BriefsAdapter, params.TeamUpdates)
			teamItemKey := params.TeamItem.Key()
			if err = tx.Delete(ctx, teamItemKey); err != nil {
				return fmt.Errorf("failed to delete team item record by key=%v: %w", teamItemKey, err)
			}
			return err
		})
}

func deleteBrief[D TeamModuleData](team record.DataWithID[string, D], itemID string, adapter *BriefsAdapter[D], updates []dal.Update) []dal.Update {
	if adapter == nil {
		return updates
	}
	count := adapter.GetBriefsCount(team.Data)
	briefIndex := -1
	for i := 0; i < count; i++ {
		if id := adapter.GetBriefItemID(team.Data, i); id == itemID {
			briefIndex = i
			break
		}
	}
	if briefIndex > 0 { // remove brief
		adapter.ShiftBriefs(team.Data, // shift all elements after found item to the left
			SliceIndexes{Start: briefIndex, End: count - 1}, // destination
			SliceIndexes{Start: briefIndex + 1, End: count}, // source
		)
		// trim list to account for the shift
		adapter.TrimBriefs(team.Data, count-1)

		// update record
		var value interface{}
		if count > 0 {
			value = adapter.BriefsValue
		}
		updates = append(updates,
			dal.Update{Field: adapter.BriefsFieldName, Value: value},
			// No need to specify counter update, created by the `facade4teamus.DeleteTeamItem`
		)
	}
	return updates
}
