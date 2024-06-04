package facade4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

//var fsRunTransaction = db.RunTransaction

var txCreate = func(ctx context.Context, inserter dal.ReadwriteTransaction, record dal.Record) error {
	return inserter.Insert(ctx, record)
}

var txUpdate = func(ctx context.Context, updater dal.ReadwriteTransaction, key *dal.Key, updates []dal.Update, preconditions ...dal.Precondition) error {
	return updater.Update(ctx, key, updates, preconditions...)
}

var txUpdateTeam = dal4teamus.TxUpdateTeam

var txCreateRetrospective = func(ctx context.Context, tx dal.ReadwriteTransaction, key *dal.Key, retrospective *dbo4retrospectus.Retrospective) error {
	if retrospective.Version == 0 {
		retrospective.Version = 1
	}
	if err := retrospective.Validate(); err != nil {
		return err
	}
	record := dal.NewRecordWithData(key, retrospective)
	return txCreate(ctx, tx, record)
}

var txUpdateRetrospective = func(ctx context.Context, tx dal.ReadwriteTransaction, key *dal.Key, retrospective *dbo4retrospectus.Retrospective, updates []dal.Update, opts ...dal.Precondition) error {
	retrospective.Version++
	return txUpdate(ctx, tx, key, append(updates, dal.Update{Field: "v", Value: retrospective.Version}), opts...)
}
