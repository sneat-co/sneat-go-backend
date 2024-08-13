package gaedal

import (
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/dal4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"time"

	"context"
)

type ReceiptDalGae struct {
}

func NewReceiptDalGae() ReceiptDalGae {
	return ReceiptDalGae{}
}

var _ dtdal.ReceiptDal = (*ReceiptDalGae)(nil)

func (ReceiptDalGae) UpdateReceipt(c context.Context, tx dal.ReadwriteTransaction, receipt models4debtus.ReceiptEntry) error {
	return tx.Set(c, receipt.Record)
}

func (receiptDalGae ReceiptDalGae) GetReceiptByID(c context.Context, tx dal.ReadSession, id string) (receipt models4debtus.ReceiptEntry, err error) {
	receipt = models4debtus.NewReceipt(id, nil)
	return receipt, tx.Get(c, receipt.Record)
}

func (receiptDalGae ReceiptDalGae) CreateReceipt(c context.Context, data *models4debtus.ReceiptDbo) (receipt models4debtus.ReceiptEntry, err error) { // TODO: Move to facade4debtus
	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		receipt = models4debtus.NewReceiptWithoutID(data)
		debtusUser := models4debtus.NewDebtusUserEntry(data.CreatorUserID)
		if err = dal4debtus.GetDebtusUser(c, tx, debtusUser); err != nil {
			return err
		}
		debtusUser.Data.CountOfReceiptsCreated += 1
		if err = tx.Set(c, debtusUser.Record); err != nil {
			return err
		}
		if err = tx.Insert(c, receipt.Record); err != nil {
			return err
		}
		receipt.ID = receipt.Record.Key().ID.(string)
		return
	})
	return
}

func (receiptDalGae ReceiptDalGae) MarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) error {
	return errors.New("TODO: Implement MarkReceiptAsSent")
	//return dtdal.DB.RunInTransaction(c, func(c context.Context) (err error) {
	//	var (
	//		receipt     models.ReceiptEntry
	//		transfer    models.TransferEntry
	//		transferKey *datastore.Key
	//	)
	//	receiptKey := NewReceiptKey(c, receiptID)
	//	if transferID == 0 {
	//		if receipt, err = receiptDalGae.GetReceiptByID(c, receiptID); err != nil {
	//			return err
	//		}
	//		if transfer, err = facade4debtus.QueueTransfers.GetTransferByID(c, transferID); err != nil {
	//			return err
	//		}
	//		transferKey = NewTransferKey(c, transferID)
	//	} else {
	//		receipt.ReceiptDbo = new(models.ReceiptDbo)
	//		transfer.TransferEntity = new(models.TransferData)
	//		transferKey = NewTransferKey(c, transferID)
	//		keys := []*datastore.Key{receiptKey, transferKey}
	//		if err = gaedb.GetMulti(c, keys, []interface{}{receipt.ReceiptDbo, transfer.TransferEntity}); err != nil {
	//			return err
	//		}
	//	}
	//
	//	if receipt.DtSent.IsZero() {
	//		receipt.DtSent = sentTime
	//		isReceiptIdIsInTransfer := false
	//		for _, rId := range transfer.ReceiptIDs {
	//			if rId == receiptID {
	//				isReceiptIdIsInTransfer = true
	//				break
	//			}
	//		}
	//		if isReceiptIdIsInTransfer {
	//			_, err = gaedb.Put(c, receiptKey, receipt)
	//		} else {
	//			transfer.ReceiptIDs = append(transfer.ReceiptIDs, receiptID)
	//			transfer.ReceiptsSentCount += 1
	//			_, err = gaedb.PutMulti(c, []*datastore.Key{receiptKey, transferKey}, []interface{}{receipt.ReceiptDbo, transfer.TransferEntity})
	//		}
	//	}
	//	return err
	//}, dtdal.CrossGroupTransaction)
}

func (receiptDalGae ReceiptDalGae) DelayedMarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) error {
	return delayerMarkReceiptAsSent.EnqueueWork(c, delaying.With(const4debtus.QueueTransfers, "set-receipt-as-sent", 0), receiptID, transferID, sentTime)
}

func delayedMarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) (err error) {
	logus.Debugf(c, "delayerMarkReceiptAsSent(receiptID=%v, transferID=%v, sentTime=%v)", receiptID, transferID, sentTime)
	if receiptID == "" {
		logus.Errorf(c, "receiptID == 0")
		return nil
	}
	if receiptID == "" {
		logus.Errorf(c, "transferID == 0")
		return nil
	}

	if err = dtdal.Receipt.MarkReceiptAsSent(c, receiptID, transferID, sentTime); dal.IsNotFound(err) {
		logus.Errorf(c, err.Error())
		return nil
	}
	return
}
