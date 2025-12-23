package facade4retrospectus

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
)

//var fsRunTransaction = db.RunTransaction

var txCreate = func(ctx context.Context, inserter dal.ReadwriteTransaction, record dal.Record) error {
	return inserter.Insert(ctx, record)
}

var txUpdate = func(ctx context.Context, updater dal.ReadwriteTransaction, key *dal.Key, updates []update.Update, preconditions ...dal.Precondition) error {
	return updater.Update(ctx, key, updates, preconditions...)
}

var txUpdateSpace = dal4spaceus.TxUpdateSpace

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

var txUpdateRetrospective = func(ctx context.Context, tx dal.ReadwriteTransaction, key *dal.Key, retrospective *dbo4retrospectus.Retrospective, updates []update.Update, opts ...dal.Precondition) error {
	retrospective.Version++
	return txUpdate(ctx, tx, key, append(updates, update.ByFieldName("v", retrospective.Version)), opts...)
}
