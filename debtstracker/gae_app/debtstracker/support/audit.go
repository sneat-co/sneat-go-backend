package support

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"reflect"
	"time"

	"context"
)

const AuditKind = "Audit"

type AuditData struct {
	Action  string
	Created time.Time
	Message string `datastore:",noindex"`
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
	c context.Context
}

func NewAuditGaeStore(c context.Context) AuditGaeStore {
	return AuditGaeStore{c: c}
}

func (s AuditGaeStore) LogAuditRecord(c context.Context, action, message string, related ...string) (audit Audit, err error) {
	audit.Data = NewAuditData(action, message, related...)
	audit.Record = dal.NewRecordWithIncompleteKey("Audit", reflect.Int, audit.Data)
	audit.Key = audit.Record.Key()
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Insert(c, audit.Record)
	})
	audit.ID = audit.Record.Key().ID.(int64)
	return
}
