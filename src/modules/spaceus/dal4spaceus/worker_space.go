package dal4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"slices"
	"strings"
	"time"
)

type spaceWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *SpaceWorkerParams) (err error)

func NewSpaceWorkerParams(userCtx facade.UserContext, spaceID string) *SpaceWorkerParams {
	return &SpaceWorkerParams{
		UserCtx: userCtx,
		Space:   dbo4spaceus.NewSpaceEntry(spaceID),
		Started: time.Now(),
	}
}

// SpaceWorkerParams passes data to a space worker
type SpaceWorkerParams struct {
	UserCtx facade.UserContext
	Started time.Time
	//
	Space         dbo4spaceus.SpaceEntry
	SpaceUpdates  []dal.Update
	RecordUpdates []RecordUpdates
}

func (v SpaceWorkerParams) UserID() string {
	return v.UserCtx.GetUserID()
}

// GetRecords loads records from DB. If userID is passed, it will check for user permissions
func (v SpaceWorkerParams) GetRecords(ctx context.Context, tx dal.ReadSession, records ...dal.Record) error {
	records = append(records, v.Space.Record)
	err := tx.GetMulti(ctx, records)
	if err != nil {
		return err
	}
	if v.UserID() != "" && v.Space.Record.Exists() {
		if !slices.Contains(v.Space.Data.UserIDs, v.UserID()) {
			return fmt.Errorf("%w: space record has no current user ContactID in UserIDs field: %s", facade.ErrUnauthorized, v.UserID())
		}
	}
	return nil
}

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
	spaceWorkerParams := NewSpaceWorkerParams(userCtx, spaceID)
	params := NewSpaceModuleWorkerParams(moduleID, spaceWorkerParams, data)
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = runModuleSpaceWorkerReadwriteTx(ctx, tx, params, worker); err != nil {
			return fmt.Errorf("failed to execute module space worker: %w", err)
		}
		return
	})
}

// RunSpaceWorkerWithUserContext executes a space worker
var RunSpaceWorkerWithUserContext = func(ctx context.Context, userCtx facade.UserContext, spaceID string, worker spaceWorker) (err error) {
	if strings.TrimSpace(spaceID) == "" {
		return fmt.Errorf("required parameter `spaceID` of RunSpaceWorkerWithUserContext() is an empty string")
	}
	if userCtx == nil {
		panic("userCtx is nil")
	}
	if userCtx.GetUserID() == "" {
		err = facade.ErrUnauthorized
		return
	}
	return runSpaceWorker(ctx, userCtx, spaceID, worker)
}

// RunSpaceWorkerWithoutUserContext executes a space worker without user context
var RunSpaceWorkerWithoutUserContext = func(ctx context.Context, spaceID string, worker spaceWorker) (err error) {
	if strings.TrimSpace(spaceID) == "" {
		return fmt.Errorf("required parameter `spaceID` of RunSpaceWorkerWithoutUserContext() is an empty string")
	}
	return runSpaceWorker(ctx, nil, spaceID, worker)
}

var runSpaceWorker = func(ctx context.Context, userCtx facade.UserContext, spaceID string, worker spaceWorker) (err error) {
	if strings.TrimSpace(spaceID) == "" {
		return fmt.Errorf("required parameter `spaceID` of runSpaceWorker() is an empty string")
	}
	userID := userCtx.GetUserID()
	if userID == "" {
		err = facade.ErrUnauthorized
		return
	}
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return RunSpaceWorkerTx(ctx, tx, userCtx, spaceID, worker)
	})
}

func RunSpaceWorkerTx(ctx context.Context, tx dal.ReadwriteTransaction, userCtx facade.UserContext, spaceID string, worker spaceWorker) (err error) {
	params := NewSpaceWorkerParams(userCtx, spaceID)
	if err = tx.Get(ctx, params.Space.Record); err != nil {
		return fmt.Errorf("failed to load space record: %w", err)
	}
	if err = params.Space.Data.Validate(); err != nil {
		return fmt.Errorf("space record loaded from DB is not valid: %w", err)
	}
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute space worker: %w", err)
	}
	if err = applySpaceUpdates(ctx, tx, params); err != nil {
		return fmt.Errorf("space worker failed to apply space record updates: %w", err)
	}
	for _, rec := range params.RecordUpdates {
		if err = tx.Update(ctx, rec.Key, rec.Updates); err != nil {
			updateFields := make([]string, len(rec.Updates))
			for _, u := range rec.Updates {
				updateField := u.Field
				if updateField == "" {
					updateField = strings.Join(u.FieldPath, ".")
				}
				updateFields = append(updateFields, updateField)
			}
			return fmt.Errorf(
				"failed to apply record updates (key=%s, updateFields: %s): %w",
				rec.Key, strings.Join(updateFields, ","), err)
		}
	}
	return
}

func applySpaceUpdates(ctx context.Context, tx dal.ReadwriteTransaction, params *SpaceWorkerParams) (err error) {
	if len(params.SpaceUpdates) == 0 && !params.Space.Record.HasChanged() {
		return
	}
	if spaceRecErr := params.Space.Record.Error(); spaceRecErr != nil {
		return fmt.Errorf("an attempt to update a space record with an error: %w", spaceRecErr)
	}
	if !params.Space.Record.HasChanged() {
		return fmt.Errorf("space record should be marked as changed before applying updates")
	}
	if err = params.Space.Data.Validate(); err != nil {
		return fmt.Errorf("space record is not valid before applying %d space updates: %w", len(params.SpaceUpdates), err)
	}
	if !params.Space.Record.Exists() {
		return tx.Insert(ctx, params.Space.Record)
	} else if len(params.SpaceUpdates) == 0 {
		return tx.Set(ctx, params.Space.Record)
	} else if err = TxUpdateSpace(ctx, tx, params.Started, params.Space, params.SpaceUpdates); err != nil {
		return fmt.Errorf("failed to update space record: %w", err)
	}
	return
}

func applySpaceModuleUpdates[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleSpaceWorkerParams[D],
) (err error) {
	if len(params.SpaceModuleUpdates) == 0 && !params.SpaceModuleEntry.Record.HasChanged() {
		return
	}
	if err = params.SpaceModuleEntry.Record.Error(); err != nil {
		return fmt.Errorf("an attempt to update a space module record that has an error: %w", err)
	}
	if !params.SpaceModuleEntry.Record.HasChanged() {
		return fmt.Errorf("space module record should be marked as changed before applying updates")
	}
	if err = params.SpaceModuleEntry.Data.Validate(); err != nil {
		return fmt.Errorf("space module record is not valid before applying space module updates: %w", err)
	}

	if params.SpaceModuleEntry.Record.Exists() {
		if len(params.SpaceModuleUpdates) == 0 {
			return tx.Set(ctx, params.SpaceModuleEntry.Record)
		} else if err = txUpdateSpaceModule(ctx, tx, params.Started, params.SpaceModuleEntry, params.SpaceModuleUpdates); err != nil {
			return fmt.Errorf("failed to update space module record: %w", err)
		}
	} else {
		if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
			return fmt.Errorf("failed to insert space module record: %w", err)
		}
	}
	return
}

// CreateSpaceItem creates a space item
func CreateSpaceItem[D SpaceModuleDbo](
	ctx context.Context,
	userCtx facade.UserContext,
	spaceRequest dto4spaceus.SpaceRequest,
	moduleID string,
	data D,
	worker func(
		ctx context.Context,
		tx dal.ReadwriteTransaction,
		spaceWorkerParams *ModuleSpaceWorkerParams[D],
	) (err error),
) (err error) {
	if worker == nil {
		panic("worker is nil")
	}
	if err := spaceRequest.Validate(); err != nil {
		return err
	}
	err = RunModuleSpaceWorker(ctx, userCtx, spaceRequest.SpaceID, moduleID, data,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *ModuleSpaceWorkerParams[D]) (err error) {
			if err := worker(ctx, tx, params); err != nil {
				return fmt.Errorf("failed to execute space worker passed to CreateSpaceItem: %w", err)
			}
			//if counter != "" {
			//	if err = incrementCounter(params.SpaceWorkerParams, moduleID, counter); err != nil {
			//		return fmt.Errorf("failed to increment spaces counter=%s: %w", counter, err)
			//	}
			//}
			if err = params.Space.Data.Validate(); err != nil {
				return fmt.Errorf("space record is not valid after performing worker: %w", err)
			}
			return
		})
	if err != nil {
		return fmt.Errorf("failed to create a space item: %w", err)
	}
	return nil
}
