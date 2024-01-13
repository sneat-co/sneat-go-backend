package facade4retrospectus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
)

//func newRetrospectiveRecord(key *dal.Key) (retrospective *dbretro.Retrospective, retrospectiveRecord dal.Record) {
//	retrospective = new(dbretro.Retrospective)
//	return retrospective, dal.NewRecordWithData(key, retrospective)
//}

var txGetRetrospective = func(ctx context.Context, tx dal.ReadwriteTransaction, record dal.Record) (err error) {
	return getRetrospective(ctx, tx, record)
}

var getRetrospective = func(ctx context.Context, getter dal.ReadSession, record dal.Record) error {
	return getter.Get(ctx, record)
}
