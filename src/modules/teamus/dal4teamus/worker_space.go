package dal4teamus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/slice"
	"strings"
	"time"
)

type spaceWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *SpaceWorkerParams) (err error)

func NewSpaceWorkerParams(userID, teamID string) *SpaceWorkerParams {
	return &SpaceWorkerParams{
		UserID:  userID,
		Space:   NewSpaceEntry(teamID),
		Started: time.Now(),
	}
}

// SpaceWorkerParams passes data to a team worker
type SpaceWorkerParams struct {
	UserID  string
	Started time.Time
	//
	Space         SpaceEntry
	SpaceUpdates  []dal.Update
	RecordUpdates []RecordUpdates
}

// GetRecords loads records from DB. If userID is passed, it will check for user permissions
func (v SpaceWorkerParams) GetRecords(ctx context.Context, tx dal.ReadSession, records ...dal.Record) error {
	records = append(records, v.Space.Record)
	err := tx.GetMulti(ctx, records)
	if err != nil {
		return err
	}
	if v.UserID != "" && v.Space.Record.Exists() {
		if !slice.Contains(v.Space.Data.UserIDs, v.UserID) {
			return fmt.Errorf("%w: team record has no current user ID in UserIDs field: %s", facade.ErrUnauthorized, v.UserID)
		}
	}
	return nil
}

// ModuleSpaceWorkerParams passes data to a team worker
type ModuleSpaceWorkerParams[D SpaceModuleDbo] struct {
	*SpaceWorkerParams
	SpaceModuleEntry   record.DataWithID[string, D]
	SpaceModuleUpdates []dal.Update
}

func (v ModuleSpaceWorkerParams[D]) GetRecords(ctx context.Context, tx dal.ReadSession, records ...dal.Record) error {
	return v.SpaceWorkerParams.GetRecords(ctx, tx, append(records, v.SpaceModuleEntry.Record)...)
}

type ModuleDbo interface {
	Validate() error
}

type SpaceModuleDbo = ModuleDbo

func RunModuleSpaceWorkerTx[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	user facade.User,
	request dto4teamus.SpaceRequest,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if worker == nil {
		panic("worker is nil")
	}
	spaceWorkerParams := NewSpaceWorkerParams(user.GetID(), request.SpaceID)
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
			dal.NewKeyWithParentAndID(spaceWorkerParams.Space.Key, SpaceModulesCollection, moduleID),
			data,
		),
	}
}

func runModuleSpaceWorkerReadonlyTx[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleSpaceWorkerParams[D],
	worker func(ctx context.Context, tx dal.ReadTransaction, teamWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute module team worker: %w", err)
	}
	return nil
}

func runModuleSpaceWorkerReadwriteTx[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleSpaceWorkerParams[D],
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute module team worker: %w", err)
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
	user facade.User,
	request dto4teamus.SpaceRequest,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadTransaction, teamWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	spaceWorkerParams := NewSpaceWorkerParams(user.GetID(), request.SpaceID)
	params := NewSpaceModuleWorkerParams(moduleID, spaceWorkerParams, data)

	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return runModuleSpaceWorkerReadonlyTx(ctx, tx, params, worker)
	})
}

func RunModuleSpaceWorker[D SpaceModuleDbo](
	ctx context.Context,
	user facade.User,
	request dto4teamus.SpaceRequest,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	teamWorkerParams := NewSpaceWorkerParams(user.GetID(), request.SpaceID)
	params := NewSpaceModuleWorkerParams(moduleID, teamWorkerParams, data)

	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return runModuleSpaceWorkerReadwriteTx(ctx, tx, params, worker)
	})
}

// RunSpaceWorker executes a team worker
var RunSpaceWorker = func(ctx context.Context, user facade.User, teamID string, worker spaceWorker) (err error) {
	if user == nil {
		panic("user is nil")
	}
	if strings.TrimSpace(teamID) == "" {
		return fmt.Errorf("spaceID is empty")
	}
	userID := user.GetID()
	if userID == "" {
		err = facade.ErrUnauthorized
		return
	}
	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		params := NewSpaceWorkerParams(userID, teamID)
		if err = tx.Get(ctx, params.Space.Record); err != nil {
			return fmt.Errorf("failed to load team record: %w", err)
		}
		if err = params.Space.Data.Validate(); err != nil {
			logus.Warningf(ctx, "Space record loaded from DB is not valid: %v: dto=%+v", err, params.Space.Data)
		}
		if err = worker(ctx, tx, params); err != nil {
			return fmt.Errorf("failed to execute team worker: %w", err)
		}
		if err = applySpaceUpdates(ctx, tx, params); err != nil {
			return fmt.Errorf("space worker failed to apply team record updates: %w", err)
		}
		for _, record := range params.RecordUpdates {
			if err = tx.Update(ctx, record.Key, record.Updates); err != nil {
				return fmt.Errorf("failed to update record (key=%s): %w", record.Key, err)
			}
		}
		return err
	})
}

func applySpaceUpdates(ctx context.Context, tx dal.ReadwriteTransaction, params *SpaceWorkerParams) (err error) {
	if len(params.SpaceUpdates) == 0 && !params.Space.Record.HasChanged() {
		return
	}
	if teamRecErr := params.Space.Record.Error(); teamRecErr != nil {
		return fmt.Errorf("an attempt to update a team record with an error: %w", teamRecErr)
	}
	if !params.Space.Record.HasChanged() {
		return fmt.Errorf("space record should be marked as changed before applying updates")
	}
	if err = params.Space.Data.Validate(); err != nil {
		return fmt.Errorf("space record is not valid before applying %d team updates: %w", len(params.SpaceUpdates), err)
	}
	if !params.Space.Record.Exists() {
		return tx.Insert(ctx, params.Space.Record)
	} else if len(params.SpaceUpdates) == 0 {
		return tx.Set(ctx, params.Space.Record)
	} else if err = TxUpdateSpace(ctx, tx, params.Started, params.Space, params.SpaceUpdates); err != nil {
		return fmt.Errorf("failed to update team record: %w", err)
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
		return fmt.Errorf("an attempt to update a team module record that has an error: %w", err)
	}
	if !params.SpaceModuleEntry.Record.HasChanged() {
		return fmt.Errorf("space module record should be marked as changed before applying updates")
	}
	if err = params.SpaceModuleEntry.Data.Validate(); err != nil {
		return fmt.Errorf("space module record is not valid before applying team module updates: %w", err)
	}

	if params.SpaceModuleEntry.Record.Exists() {
		if err = txUpdateSpaceModule(ctx, tx, params.Started, params.SpaceModuleEntry, params.SpaceModuleUpdates); err != nil {
			return fmt.Errorf("failed to update team module record: %w", err)
		}
	} else {
		if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
			return fmt.Errorf("failed to insert team module record: %w", err)
		}
	}
	return
}

// CreateSpaceItem creates a team item
func CreateSpaceItem[D SpaceModuleDbo](
	ctx context.Context,
	user facade.User,
	teamRequest dto4teamus.SpaceRequest,
	moduleID string,
	data D,
	worker func(
		ctx context.Context,
		tx dal.ReadwriteTransaction,
		teamWorkerParams *ModuleSpaceWorkerParams[D],
	) (err error),
) (err error) {
	if worker == nil {
		panic("worker is nil")
	}
	if err := teamRequest.Validate(); err != nil {
		return err
	}
	err = RunModuleSpaceWorker(ctx, user, teamRequest, moduleID, data,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *ModuleSpaceWorkerParams[D]) (err error) {
			if err := worker(ctx, tx, params); err != nil {
				return fmt.Errorf("failed to execute team worker passed to CreateSpaceItem: %w", err)
			}
			//if counter != "" {
			//	if err = incrementCounter(params.SpaceWorkerParams, moduleID, counter); err != nil {
			//		return fmt.Errorf("failed to incement teams counter=%s: %w", counter, err)
			//	}
			//}
			if err = params.Space.Data.Validate(); err != nil {
				return fmt.Errorf("space record is not valid after performing worker: %w", err)
			}
			return
		})
	if err != nil {
		return fmt.Errorf("failed to create a team item: %w", err)
	}
	return nil
}
