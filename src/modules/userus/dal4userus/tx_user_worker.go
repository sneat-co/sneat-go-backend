package dal4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
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

var RunUserWorker = func(ctx context.Context, userCtx facade.UserContext, worker userWorker) (err error) {
	if userCtx == nil {
		panic("userCtx == nil")
	}
	params := UserWorkerParams{
		User:    dbo4userus.NewUserEntry(userCtx.GetUserID()),
		Started: time.Now(),
	}
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = tx.Get(ctx, params.User.Record); err != nil {
			return fmt.Errorf("failed to load userCtx record: %w", err)
		}
		if err = params.User.Data.Validate(); err != nil {
			err = fmt.Errorf("user record loaded from DB is not valid: %v: userID=%s data=%+v", err, params.User.ID, params.User.Data)
			return
		}
		if err = worker(ctx, tx, &params); err != nil {
			return fmt.Errorf("failed to execute teamWorker: %w", err)
		}
		if err = params.User.Data.Validate(); err != nil {
			return fmt.Errorf("user record is not valid after userWorker completion: %w", err)
		}
		if len(params.UserUpdates) > 0 {
			if err = TxUpdateUser(ctx, tx, params.Started, params.User.Key, params.UserUpdates); err != nil {
				return fmt.Errorf("failed to update team record: %w", err)
			}
		}
		return err
	})
	if err != nil {
		err = fmt.Errorf("failed inside transaction created by RunUserWorker(): %w", err)
	}
	return
}
