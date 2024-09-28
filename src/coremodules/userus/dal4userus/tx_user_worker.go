package dal4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"time"
)

// UserWorkerParams passes data to a team worker
type UserWorkerParams struct {
	Started     time.Time
	User        dbo4userus.UserEntry
	UserUpdates []dal.Update
}

type userWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, userWorkerParams *UserWorkerParams) (err error)

var RunUserWorker = func(ctx context.Context, userCtx facade.UserContext, userRecordMustExists bool, worker userWorker) (err error) {
	if userCtx == nil {
		panic("userCtx == nil")
	}
	params := UserWorkerParams{
		User:    dbo4userus.NewUserEntry(userCtx.GetUserID()),
		Started: time.Now(),
	}
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = tx.Get(ctx, params.User.Record); err != nil {
			if dal.IsNotFound(err) && !userRecordMustExists {
				params.User.Record.SetError(dal.ErrRecordNotFound)
			} else {
				return fmt.Errorf("failed to load user record: %w", err)
			}
		} else if err = params.User.Data.Validate(); err != nil { // Do not validate if user record is not found
			err = fmt.Errorf("user record loaded from DB is not valid: %v: userID=%s data=%+v", err, params.User.ID, params.User.Data)
			return
		}
		if err = worker(ctx, tx, &params); err != nil {
			return fmt.Errorf("failed to execute teamWorker: %w", err)
		}
		if err = params.User.Data.Validate(); err != nil {
			return fmt.Errorf("user record is not valid after userWorker completion: %w", err)
		}
		if params.User.Record.HasChanged() {
			if params.User.Record.Exists() {
				if len(params.UserUpdates) > 0 {
					if err = TxUpdateUser(ctx, tx, params.Started, params.User.Key, params.UserUpdates); err != nil {
						return fmt.Errorf("failed to apply updates to user record: %w", err)
					}
				} else {
					if err = tx.Set(ctx, params.User.Record); err != nil {
						return fmt.Errorf("failed to update user record: %w", err)
					}
				}
			} else {
				if err = tx.Insert(ctx, params.User.Record); err != nil {
					return fmt.Errorf("failed to insert user record: %w", err)
				}
			}
		} else if len(params.UserUpdates) > 0 {
			return fmt.Errorf("user record is not marked as changed but there are updates to apply")
		}
		return err
	})
	if err != nil {
		err = fmt.Errorf("failed inside transaction created by RunUserWorker(): %w", err)
	}
	return
}
