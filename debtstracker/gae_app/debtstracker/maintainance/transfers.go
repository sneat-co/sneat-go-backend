package maintainance

//import (
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
//	"context"
//	"github.com/captaincodeman/datastore-mapper"
//	"github.com/dal-go/dalgo/dal"
//	"net/http"
//)
//
//type transfersAsyncJob struct {
//	asyncMapper
//	entity *models.TransferData
//}
//
//func (m *transfersAsyncJob) Make() interface{} {
//	m.entity = new(models.TransferData)
//	return m.entity
//}
//
//func (m *transfersAsyncJob) Query(r *http.Request) (query *mapper.Query, err error) {
//	return applyIDAndUserFilters(r, "transfersAsyncJob", models.TransferKind, filterByIntID, "BothUserIDs")
//}
//
//func (m *transfersAsyncJob) Transfer(key *dal.Key) models.Transfer {
//	entity := *m.entity
//	return models.NewTransfer(key.ID.(int), &entity)
//}
//
//type TransferWorker func(c context.Context, tx dal.ReadwriteTransaction, counters *asyncCounters, transfer models.Transfer) error

//func (m *transfersAsyncJob) startTransferWorker(c context.Context, counters mapper.Counters, key *dal.Key, transferWorker TransferWorker) error {
//	transfer := m.Transfer(key)
//	w := func() Worker {
//		return func(counters *asyncCounters) error {
//			db, err := facade.GetDatabase(c)
//			if err != nil {
//				return err
//			}
//			return db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
//				return transferWorker(c, tx, counters, transfer)
//			})
//
//		}
//	}
//	return m.startWorker(c, counters, w)
//}
