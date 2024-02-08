package gaedal

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
)

func getDatabase(_ context.Context) (db dal.DB, err error) {
	return nil, errors.New("TODO: implement me: GetDatabase()")
}

var GetDatabase func(context.Context) (dal.DB, error) = getDatabase
