package gaedal

import (
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/delaying"
	"time"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

type ReceiptDalGae struct {
}

func NewReceiptDalGae() ReceiptDalGae {
	return ReceiptDalGae{}
}

var _ dtdal.ReceiptDal = (*ReceiptDalGae)(nil)

func (ReceiptDalGae) UpdateReceipt(c context.Context, tx dal.ReadwriteTransaction, receipt models.Receipt) error {
	return tx.Set(c, receipt.Record)
}

func (receiptDalGae ReceiptDalGae) GetReceiptByID(c context.Context, tx dal.ReadSession, id string) (receipt models.Receipt, err error) {
	receipt = models.NewReceipt(id, nil)
	return receipt, tx.Get(c, receipt.Record)
}

func (receiptDalGae ReceiptDalGae) CreateReceipt(c context.Context, data *models.ReceiptData) (receipt models.Receipt, err error) { // TODO: Move to facade
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		receipt = models.NewReceiptWithoutID(data)
		var user models.AppUser
		if user, err = facade.User.GetUserByID(c, tx, data.CreatorUserID); err != nil {
			return err
		}
		user.Data.CountOfReceiptsCreated += 1
		if err = tx.Set(c, user.Record); err != nil {
			return err
		}
		if err = tx.Insert(c, receipt.Record); err != nil {
			return err
		}
		receipt.ID = receipt.Record.Key().ID.(string)
		return nil
	})
	return
}

func (receiptDalGae ReceiptDalGae) MarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) error {
	return errors.New("TODO: Implement MarkReceiptAsSent")
	//return dtdal.DB.RunInTransaction(c, func(c context.Context) (err error) {
	//	var (
	//		receipt     models.Receipt
	//		transfer    models.Transfer
	//		transferKey *datastore.Key
	//	)
	//	receiptKey := NewReceiptKey(c, receiptID)
	//	if transferID == 0 {
	//		if receipt, err = receiptDalGae.GetReceiptByID(c, receiptID); err != nil {
	//			return err
	//		}
	//		if transfer, err = facade.Transfers.GetTransferByID(c, transferID); err != nil {
	//			return err
	//		}
	//		transferKey = NewTransferKey(c, transferID)
	//	} else {
	//		receipt.ReceiptData = new(models.ReceiptData)
	//		transfer.TransferEntity = new(models.TransferData)
	//		transferKey = NewTransferKey(c, transferID)
	//		keys := []*datastore.Key{receiptKey, transferKey}
	//		if err = gaedb.GetMulti(c, keys, []interface{}{receipt.ReceiptData, transfer.TransferEntity}); err != nil {
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
	//			_, err = gaedb.PutMulti(c, []*datastore.Key{receiptKey, transferKey}, []interface{}{receipt.ReceiptData, transfer.TransferEntity})
	//		}
	//	}
	//	return err
	//}, dtdal.CrossGroupTransaction)
}

func (receiptDalGae ReceiptDalGae) DelayedMarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) error {
	return delayMarkReceiptAsSent.EnqueueWork(c, delaying.With(common.QUEUE_TRANSFERS, "set-receipt-as-sent", 0), receiptID, transferID, sentTime)
}

func delayedMarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) (err error) {
	log.Debugf(c, "delayMarkReceiptAsSent(receiptID=%v, transferID=%v, sentTime=%v)", receiptID, transferID, sentTime)
	if receiptID == "" {
		log.Errorf(c, "receiptID == 0")
		return nil
	}
	if receiptID == "" {
		log.Errorf(c, "transferID == 0")
		return nil
	}

	if err = dtdal.Receipt.MarkReceiptAsSent(c, receiptID, transferID, sentTime); dal.IsNotFound(err) {
		log.Errorf(c, err.Error())
		return nil
	}
	return
}
