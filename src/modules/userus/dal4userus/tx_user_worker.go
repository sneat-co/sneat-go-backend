package dal4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"time"
)

// UserWorkerParams passes data to a team worker
type UserWorkerParams struct {
	Started     time.Time
	User        dbo4userus.User
	UserUpdates []dal.Update
}

type userWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, userWorkerParams *UserWorkerParams) (err error)

var RunUserWorker = func(ctx context.Context, userCtx facade.UserContext, worker userWorker) (err error) {
	if userCtx == nil {
		panic("userCtx == nil")
	}
	params := UserWorkerParams{
		User:    dbo4userus.NewUser(userCtx.GetUserID()),
		Started: time.Now(),
	}
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = tx.Get(ctx, params.User.Record); err != nil {
			return fmt.Errorf("failed to load userCtx record: %w", err)
		}
		if err = params.User.Data.Validate(); err != nil {
			logus.Warningf(ctx, "User record loaded from DB is not valid: %v: data=%+v", err, params.User.Data)
		}
		if err = worker(ctx, tx, &params); err != nil {
			return fmt.Errorf("failed to execute teamWorker: %w", err)
		}
		if err = params.User.Data.Validate(); err != nil {
			return fmt.Errorf("userCtx record is not valid after update: %w", err)
		}
		if len(params.UserUpdates) > 0 {
			if err = TxUpdateUser(ctx, tx, params.Started, params.User.Key, params.UserUpdates); err != nil {
				return fmt.Errorf("failed to update team record: %w", err)
			}
		}
		return err
	})
}
