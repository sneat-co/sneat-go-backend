package maintainance

//import (
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/facade4debtus"
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/models"
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
//		logus.Errorf(c, "Transfer(ContactID=%v) is missing CreatorUserID")
//		return
//	}
//	if !transfer.Data.HasObsoleteProps() {
//		// logus.Debugf(c, "transfer.ContactID=%v has no obsolete props", transfer.ContactID)
//		return
//	}
//	var db dal.DB
//	if db, err = facade4debtus.GetDatabase(c); err != nil {
//		return err
//	}
//
//	if err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
//		if transfer, err = facade4debtus.Transfers.GetTransferByID(c, tx, transfer.ContactID); err != nil {
//			return
//		}
//		if transfer.Data.HasObsoleteProps() {
//			if err = facade4debtus.Transfers.SaveTransfer(tc, tx, transfer); err != nil {
//				return
//			}
//			logus.Infof(c, "Transfer %v fixed", transfer.ContactID)
//		}
//		return
//	}, nil); err != nil {
//		logus.Errorf(c, "failed to fix transfer %v: %v", transfer.ContactID, err)
//	}
//	return
//}
