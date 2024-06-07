package dal4teamus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
	"log"
	"strings"
	"time"
)

type teamWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *TeamWorkerParams) (err error)

func NewTeamWorkerParams(userID, teamID string) *TeamWorkerParams {
	return &TeamWorkerParams{
		UserID:  userID,
		Team:    NewTeamEntry(teamID),
		Started: time.Now(),
	}
}

// TeamWorkerParams passes data to a team worker
type TeamWorkerParams struct {
	UserID  string
	Started time.Time
	//
	Team          TeamEntry
	TeamUpdates   []dal.Update
	RecordUpdates []RecordUpdates
}

// GetRecords loads records from DB. If userID is passed, it will check for user permissions
func (v TeamWorkerParams) GetRecords(ctx context.Context, tx dal.ReadSession, records ...dal.Record) error {
	records = append(records, v.Team.Record)
	err := tx.GetMulti(ctx, records)
	if err != nil {
		return err
	}
	if v.UserID != "" && v.Team.Record.Exists() {
		if !slice.Contains(v.Team.Data.UserIDs, v.UserID) {
			return fmt.Errorf("%w: team record has no current user ID in UserIDs field: %s", facade.ErrUnauthorized, v.UserID)
		}
	}
	return nil
}

// ModuleTeamWorkerParams passes data to a team worker
type ModuleTeamWorkerParams[D TeamModuleDbo] struct {
	*TeamWorkerParams
	TeamModuleEntry   record.DataWithID[string, D]
	TeamModuleUpdates []dal.Update
}

func (v ModuleTeamWorkerParams[D]) GetRecords(ctx context.Context, tx dal.ReadSession, records ...dal.Record) error {
	return v.TeamWorkerParams.GetRecords(ctx, tx, append(records, v.TeamModuleEntry.Record)...)
}

type ModuleDbo interface {
	Validate() error
}

type TeamModuleDbo = ModuleDbo

func RunModuleTeamWorkerTx[D TeamModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	user facade.User,
	request dto4teamus.TeamRequest,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleTeamWorkerParams[D]) (err error),
) (err error) {
	if worker == nil {
		panic("worker is nil")
	}
	teamWorkerParams := NewTeamWorkerParams(user.GetID(), request.TeamID)
	params := NewTeamModuleWorkerParams(moduleID, teamWorkerParams, data)
	return runModuleTeamWorkerReadwriteTx(ctx, tx, params, worker)
}

func NewTeamModuleWorkerParams[D TeamModuleDbo](
	moduleID string,
	teamWorkerParams *TeamWorkerParams,
	data D,
) *ModuleTeamWorkerParams[D] {
	return &ModuleTeamWorkerParams[D]{
		TeamWorkerParams: teamWorkerParams,
		TeamModuleEntry: record.NewDataWithID(moduleID,
			dal.NewKeyWithParentAndID(teamWorkerParams.Team.Key, TeamModulesCollection, moduleID),
			data,
		),
	}
}

func runModuleTeamWorkerReadonlyTx[D TeamModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleTeamWorkerParams[D],
	worker func(ctx context.Context, tx dal.ReadTransaction, teamWorkerParams *ModuleTeamWorkerParams[D]) (err error),
) (err error) {
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute module team worker: %w", err)
	}
	return nil
}

func runModuleTeamWorkerReadwriteTx[D TeamModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleTeamWorkerParams[D],
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleTeamWorkerParams[D]) (err error),
) (err error) {
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute module team worker: %w", err)
	}
	if err = applyTeamUpdates(ctx, tx, params.TeamWorkerParams); err != nil {
		return fmt.Errorf("module team worker failed to apply team record updates: %w", err)
	}
	if err = applyTeamModuleUpdates(ctx, tx, params); err != nil {
		return fmt.Errorf("module team worker failed to apply team module record updates: %w", err)
	}
	return nil
}

func RunReadonlyModuleTeamWorker[D TeamModuleDbo](
	ctx context.Context,
	user facade.User,
	request dto4teamus.TeamRequest,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadTransaction, teamWorkerParams *ModuleTeamWorkerParams[D]) (err error),
) (err error) {
	teamWorkerParams := NewTeamWorkerParams(user.GetID(), request.TeamID)
	params := NewTeamModuleWorkerParams(moduleID, teamWorkerParams, data)

	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return runModuleTeamWorkerReadonlyTx(ctx, tx, params, worker)
	})
}

func RunModuleTeamWorker[D TeamModuleDbo](
	ctx context.Context,
	user facade.User,
	request dto4teamus.TeamRequest,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleTeamWorkerParams[D]) (err error),
) (err error) {
	teamWorkerParams := NewTeamWorkerParams(user.GetID(), request.TeamID)
	params := NewTeamModuleWorkerParams(moduleID, teamWorkerParams, data)

	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return runModuleTeamWorkerReadwriteTx(ctx, tx, params, worker)
	})
}

// RunTeamWorker executes a team worker
var RunTeamWorker = func(ctx context.Context, user facade.User, teamID string, worker teamWorker) (err error) {
	if user == nil {
		panic("user is nil")
	}
	if strings.TrimSpace(teamID) == "" {
		return fmt.Errorf("teamID is empty")
	}
	userID := user.GetID()
	if userID == "" {
		err = facade.ErrUnauthorized
		return
	}
	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		params := NewTeamWorkerParams(userID, teamID)
		if err = tx.Get(ctx, params.Team.Record); err != nil {
			return fmt.Errorf("failed to load team record: %w", err)
		}
		if err = params.Team.Data.Validate(); err != nil {
			log.Printf("WARNING: team record loaded from DB is not valid: %v: dto=%+v", err, params.Team.Data)
		}
		if err = worker(ctx, tx, params); err != nil {
			return fmt.Errorf("failed to execute team worker: %w", err)
		}
		if err = applyTeamUpdates(ctx, tx, params); err != nil {
			return fmt.Errorf("team worker failed to apply team record updates: %w", err)
		}
		for _, record := range params.RecordUpdates {
			if err = tx.Update(ctx, record.Key, record.Updates); err != nil {
				return fmt.Errorf("failed to update record (key=%s): %w", record.Key, err)
			}
		}
		return err
	})
}

func applyTeamUpdates(ctx context.Context, tx dal.ReadwriteTransaction, params *TeamWorkerParams) (err error) {
	if len(params.TeamUpdates) > 0 {
		if err = params.Team.Data.Validate(); err != nil {
			return fmt.Errorf("team record is not valid before applying %d team updates: %w", len(params.TeamUpdates), err)
		}
		if err = TxUpdateTeam(ctx, tx, params.Started, params.Team, params.TeamUpdates); err != nil {
			return fmt.Errorf("failed to update team record: %w", err)
		}
	}
	return err
}

func applyTeamModuleUpdates[D TeamModuleDbo](ctx context.Context, tx dal.ReadwriteTransaction, params *ModuleTeamWorkerParams[D]) (err error) {
	if len(params.TeamModuleUpdates) > 0 {
		if err = params.TeamModuleEntry.Data.Validate(); err != nil {
			return fmt.Errorf("team module record is not valid before applying team module updates: %w", err)
		}
		if params.TeamModuleEntry.Record.Exists() {
			if err = txUpdateTeamModule(ctx, tx, params.Started, params.TeamModuleEntry, params.TeamModuleUpdates); err != nil {
				return fmt.Errorf("failed to update team module record: %w", err)
			}
		} else {
			if err = tx.Insert(ctx, params.TeamModuleEntry.Record); err != nil {
				return fmt.Errorf("failed to insert team module record: %w", err)
			}
		}
	}
	return err
}

// CreateTeamItem creates a team item
func CreateTeamItem[D TeamModuleDbo](
	ctx context.Context,
	user facade.User,
	teamRequest dto4teamus.TeamRequest,
	moduleID string,
	data D,
	worker func(
		ctx context.Context,
		tx dal.ReadwriteTransaction,
		teamWorkerParams *ModuleTeamWorkerParams[D],
	) (err error),
) (err error) {
	if worker == nil {
		panic("worker is nil")
	}
	if err := teamRequest.Validate(); err != nil {
		return err
	}
	err = RunModuleTeamWorker(ctx, user, teamRequest, moduleID, data,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *ModuleTeamWorkerParams[D]) (err error) {
			if err := worker(ctx, tx, params); err != nil {
				return fmt.Errorf("failed to execute team worker passed to CreateTeamItem: %w", err)
			}
			//if counter != "" {
			//	if err = incrementCounter(params.TeamWorkerParams, moduleID, counter); err != nil {
			//		return fmt.Errorf("failed to incement teams counter=%s: %w", counter, err)
			//	}
			//}
			if err = params.Team.Data.Validate(); err != nil {
				return fmt.Errorf("team record is not valid after performing worker: %w", err)
			}
			return
		})
	if err != nil {
		return fmt.Errorf("failed to create a team item: %w", err)
	}
	return nil
}
