package support

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-core/facade"
	"reflect"
	"time"

	"context"
)

const AuditKind = "Audit"

type AuditData struct {
	Action  string
	Created time.Time
	Message string `firestore:",omitempty"`
	Related []string
}

type Audit struct {
	record.WithID[int64]
	Data *AuditData
}

func NewAuditKey(id int64) *dal.Key {
	return dal.NewKeyWithID(AuditKind, id)
}

func NewAudit(id int64, data *AuditData) Audit {
	key := NewAuditKey(id)
	return Audit{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}

func NewAuditData(action, message string, related ...string) *AuditData {
	return &AuditData{
		Created: time.Now(),
		Action:  action,
		Message: message,
		Related: related,
	}
}

type AuditStorage interface {
	LogAuditRecord(action, message string, related ...string) error
}

type AuditGaeStore struct {
	ctx context.Context
}

func NewAuditGaeStore(ctx context.Context) AuditGaeStore {
	return AuditGaeStore{ctx: ctx}
}

func (s AuditGaeStore) LogAuditRecord(ctx context.Context, action, message string, related ...string) (audit Audit, err error) {
	audit.Data = NewAuditData(action, message, related...)
	audit.Record = dal.NewRecordWithIncompleteKey("Audit", reflect.Int, audit.Data)
	audit.Key = audit.Record.Key()
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Insert(ctx, audit.Record)
	})
	audit.ID = audit.Record.Key().ID.(int64)
	return
}
