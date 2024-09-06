package gaedal

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"

	//"errors"
	"sync"

	"context"
)

type TransferFixer struct {
	changed     bool
	Fixes       []string
	transferKey *dal.Key
	transfer    *models4debtus.TransferData
}

func NewTransferFixer(transferKey *dal.Key, transfer *models4debtus.TransferData) TransferFixer {
	return TransferFixer{transferKey: transferKey, transfer: transfer, Fixes: make([]string, 0)}
}

func (f *TransferFixer) needFixCounterpartyCounterpartyName() bool {
	return f.transfer.Creator().ContactName == ""
}

//func (f *TransferFixer) fixCounterpartyCounterpartyName(ctx context.Context) error {
//	if f.needFixCounterpartyCounterpartyName() {
//		logus.Debugf(ctx, "%v: needFixCounterpartyCounterpartyName=true", f.transferKey.IntegerID())
//		if f.transfer.Creator().CounterpartyID != 0 {
//			var counterpartyCounterparty models.DebtusSpaceContactDbo
//			err := gaedb.Get(c, NewCounterpartyKey(c, f.transfer.Creator().CounterpartyID), &counterpartyCounterparty)
//			if err != nil {
//				return err
//			}
//			f.transfer.Creator().ContactName = counterpartyCounterparty.GetFullName()
//			logus.Debugf(c, "%v: got name from counterpartyCounterparty", f.transferKey.IntegerID())
//			if f.transfer.Creator().ContactName == "" {
//				logus.Warningf(c, "Counterparty %v has no full name", f.transfer.Creator().CounterpartyID)
//			}
//		}
//		if f.transfer.Creator().ContactName == "" { // Not fixed from counterparty
//			user, err := dal4userus.GetUserByID(c, f.transfer.CreatorUserID)
//			if err != nil {
//				return err
//			}
//			f.transfer.Creator().ContactName = user.GetFullName()
//			logus.Debugf(c, "%v: got name from user", f.transferKey.IntegerID())
//			if f.transfer.Creator().ContactName == "" {
//				logus.Warningf(c, "User %v has no full name", f.transfer.CreatorUserID)
//			}
//		}
//		if f.transfer.Creator().ContactName == "" {
//			return errors.New("f.transfer.Creator().ContactName is not fixed")
//		}
//		f.changed = true
//		f.Fixes = append(f.Fixes, "CounterpartyCounterpartyName")
//		//} else {
//		//	logus.Debugf(ctx, "%v: %v", f.transferKey.IntegerID(), f.transfer.Creator().ContactName)
//	}
//	return nil
//}

func (f *TransferFixer) needFixes(_ context.Context) bool {
	return f.needFixCounterpartyCounterpartyName()
	//logus.Debugf(c, "%v: needFixes=%v", f.transferKey.IntegerID(), result)
	//return result
}

func (f *TransferFixer) FixAllIfNeeded(ctx context.Context) (err error) {
	if f.needFixes(ctx) {
		err = facade.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) error {
			transfer, err := facade4debtus.Transfers.GetTransferByID(tctx, tx, f.transferKey.ID.(string))
			if err != nil {
				return err
			}
			f.transfer = transfer.Data
			//if err = f.fixCounterpartyCounterpartyName(ctx); err != nil {
			//	return err
			//}
			if f.changed {
				//logus.Debugf(ctx, "%v: changed", f.transferKey.IntegerID())
				err = tx.Set(tctx, transfer.Record)
				return err
				//} else {
				//	logus.Debugf(ctx, "%v: not changed", f.transferKey.IntegerID())
			}
			return nil
		}, nil)
	}
	return
}

func FixTransfers(ctx context.Context) (loadedCount int, fixedCount int, failedCount int, err error) {
	query := dal.From(models4debtus.TransfersCollection).SelectInto(func() dal.Record {
		return models4debtus.NewTransferWithIncompleteKey(nil).Record
	})
	//query.Limit = 50
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	var reader dal.Reader
	reader, err = db.QueryReader(ctx, query)
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for {
		var record dal.Record
		if record, err = reader.Next(); err != nil {
			if err == dal.ErrNoMoreRecords {
				err = nil
				return
			}
			logus.Errorf(ctx, "Failed to get next transfer: %v", err.Error())
			return
		}
		loadedCount += 1
		wg.Add(1)
		go func(transferRecord dal.Record) {
			defer wg.Done()
			key := transferRecord.Key()
			fixer := NewTransferFixer(key, transferRecord.Data().(*models4debtus.TransferData))
			err2 := fixer.FixAllIfNeeded(ctx)
			if err2 != nil {
				logus.Errorf(ctx, "Failed to fix transfer=%v: %v", key.ID.(int), err2.Error())
				mutex.Lock()
				failedCount += 1
				err = err2
				mutex.Unlock()
			} else {
				if len(fixer.Fixes) > 0 {
					mutex.Lock()
					fixedCount += 1
					mutex.Unlock()
					logus.Infof(ctx, "Fixed transfer %v: %v", key.ID.(int), fixer.Fixes)
					//} else {
					//	logus.Debugf(ctx, "TransferEntry %v is OK: CounterpartyCounterpartyName: %v", transferKey.IntegerID(), fixer.transfer.Creator().ContactName)
				}
			}
		}(record)
		if err != nil {
			break
		}
	}
	wg.Wait()
	return
}
