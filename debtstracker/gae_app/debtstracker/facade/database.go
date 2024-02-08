package facade

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
)

// GetDatabase returns debts tracker database
func GetDatabase(ctx context.Context) (db dal.DB, err error) {
	return nil, errors.New("TODO: Implement GetDatabase()")
}

func DB() dal.DB {
	panic("TODO: Implement DB()")
}

func RunReadwriteTransaction(ctx context.Context, f func(ctx context.Context, tx dal.ReadwriteTransaction) error) error {
	db, err := GetDatabase(ctx)
	if err != nil {
		return err
	}
	return db.RunReadwriteTransaction(ctx, f)
}
