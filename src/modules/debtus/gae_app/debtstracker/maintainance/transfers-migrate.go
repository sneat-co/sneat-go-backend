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
//func (m *migrateTransfers) Next(ctx context.Context, counters mapper.Counters, key *dal.Key) (err error) {
//	return m.startTransferWorker(ctx, counters, key, m.migrateTransfer)
//}
//
//func (m *migrateTransfers) migrateTransfer(ctx context.Context, tx dal.ReadwriteTransaction, counters *asyncCounters, transfer models.Transfer) (err error) {
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
//	if err = db.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) (err error) {
//		if transfer, err = facade4debtus.Transfers.GetTransferByID(tctx, tx, transfer.ContactID); err != nil {
//			return
//		}
//		if transfer.Data.HasObsoleteProps() {
//			if err = facade4debtus.Transfers.SaveTransfer(tctx, tx, transfer); err != nil {
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
