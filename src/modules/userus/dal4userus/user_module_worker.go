package dal4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-core/facade"
)

type UserModuleDbo = interface {
	Validate() error
}

type UserModuleWorkerParams[T any] struct {
	UserModule        record.DataWithID[string, *T]
	UserModuleUpdates []dal.Update
}

func RunUserModuleWorker[T any](
	ctx context.Context,
	userID, moduleID string,
	userModuleDbo *T,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, param *UserModuleWorkerParams[T]) error,
) (err error) {
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		userModule := record.NewDataWithID(userID, NewUserModuleKey(userID, moduleID), userModuleDbo)
		params := UserModuleWorkerParams[T]{
			UserModule: userModule,
		}
		if err = worker(ctx, tx, &params); err != nil {
			return err
		}
		if len(params.UserModuleUpdates) > 0 {
			if err = tx.Update(ctx, params.UserModule.Key, params.UserModuleUpdates); err != nil {
				return fmt.Errorf("failed to update user module: %w", err)
			}
		}
		return err
	})
	return err
}
