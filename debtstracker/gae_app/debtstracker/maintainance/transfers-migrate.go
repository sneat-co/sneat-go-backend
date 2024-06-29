package maintainance

//import (
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
//	"context"
//	"github.com/captaincodeman/datastore-mapper"
//	"github.com/dal-go/dalgo/dal"
//	"github.com/strongo/logus"
//)
//
//type migrateTransfers struct {
//	transfersAsyncJob
//}
//
//func (m *migrateTransfers) Next(c context.Context, counters mapper.Counters, key *dal.Key) (err error) {
//	return m.startTransferWorker(c, counters, key, m.migrateTransfer)
//}
//
//func (m *migrateTransfers) migrateTransfer(c context.Context, tx dal.ReadwriteTransaction, counters *asyncCounters, transfer models.Transfer) (err error) {
//	if transfer.Data.CreatorUserID == 0 {
//		logus.Errorf(c, "Transfer(ID=%v) is missing CreatorUserID")
//		return
//	}
//	if !transfer.Data.HasObsoleteProps() {
//		// logus.Debugf(c, "transfer.ID=%v has no obsolete props", transfer.ID)
//		return
//	}
//	var db dal.DB
//	if db, err = facade.GetDatabase(c); err != nil {
//		return err
//	}
//
//	if err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
//		if transfer, err = facade.Transfers.GetTransferByID(c, tx, transfer.ID); err != nil {
//			return
//		}
//		if transfer.Data.HasObsoleteProps() {
//			if err = facade.Transfers.SaveTransfer(tc, tx, transfer); err != nil {
//				return
//			}
//			logus.Infof(c, "Transfer %v fixed", transfer.ID)
//		}
//		return
//	}, nil); err != nil {
//		logus.Errorf(c, "failed to fix transfer %v: %v", transfer.ID, err)
//	}
//	return
//}
