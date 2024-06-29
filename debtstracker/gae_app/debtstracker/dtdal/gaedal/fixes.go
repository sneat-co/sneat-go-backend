package gaedal

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/logus"

	//"errors"
	"sync"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
)

type TransferFixer struct {
	changed     bool
	Fixes       []string
	transferKey *dal.Key
	transfer    *models.TransferData
}

func NewTransferFixer(transferKey *dal.Key, transfer *models.TransferData) TransferFixer {
	return TransferFixer{transferKey: transferKey, transfer: transfer, Fixes: make([]string, 0)}
}

func (f *TransferFixer) needFixCounterpartyCounterpartyName() bool {
	return f.transfer.Creator().ContactName == ""
}

//func (f *TransferFixer) fixCounterpartyCounterpartyName(c context.Context) error {
//	if f.needFixCounterpartyCounterpartyName() {
//		logus.Debugf(c, "%v: needFixCounterpartyCounterpartyName=true", f.transferKey.IntegerID())
//		if f.transfer.Creator().CounterpartyID != 0 {
//			var counterpartyCounterparty models.DebtusContactDbo
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
//			user, err := facade.User.GetUserByID(c, f.transfer.CreatorUserID)
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
//		//	logus.Debugf(c, "%v: %v", f.transferKey.IntegerID(), f.transfer.Creator().ContactName)
//	}
//	return nil
//}

func (f *TransferFixer) needFixes(_ context.Context) bool {
	return f.needFixCounterpartyCounterpartyName()
	//logus.Debugf(c, "%v: needFixes=%v", f.transferKey.IntegerID(), result)
	//return result
}

func (f *TransferFixer) FixAllIfNeeded(c context.Context) (err error) {
	if f.needFixes(c) {
		var db dal.DB
		if db, err = facade.GetDatabase(c); err != nil {
			return
		}

		err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
			transfer, err := facade.Transfers.GetTransferByID(tc, tx, f.transferKey.ID.(string))
			if err != nil {
				return err
			}
			f.transfer = transfer.Data
			//if err = f.fixCounterpartyCounterpartyName(c); err != nil {
			//	return err
			//}
			if f.changed {
				//logus.Debugf(c, "%v: changed", f.transferKey.IntegerID())
				err = tx.Set(tc, transfer.Record)
				return err
				//} else {
				//	logus.Debugf(c, "%v: not changed", f.transferKey.IntegerID())
			}
			return nil
		}, nil)
	}
	return
}

func FixTransfers(c context.Context) (loadedCount int, fixedCount int, failedCount int, err error) {
	query := dal.From(models.TransfersCollection).SelectInto(func() dal.Record {
		return models.NewTransferWithIncompleteKey(nil).Record
	})
	//query.Limit = 50
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	var reader dal.Reader
	reader, err = db.QueryReader(c, query)
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
			logus.Errorf(c, "Failed to get next transfer: %v", err.Error())
			return
		}
		loadedCount += 1
		wg.Add(1)
		go func(transferRecord dal.Record) {
			defer wg.Done()
			key := transferRecord.Key()
			fixer := NewTransferFixer(key, transferRecord.Data().(*models.TransferData))
			err2 := fixer.FixAllIfNeeded(c)
			if err2 != nil {
				logus.Errorf(c, "Failed to fix transfer=%v: %v", key.ID.(int), err2.Error())
				mutex.Lock()
				failedCount += 1
				err = err2
				mutex.Unlock()
			} else {
				if len(fixer.Fixes) > 0 {
					mutex.Lock()
					fixedCount += 1
					mutex.Unlock()
					logus.Infof(c, "Fixed transfer %v: %v", key.ID.(int), fixer.Fixes)
					//} else {
					//	logus.Debugf(c, "TransferEntry %v is OK: CounterpartyCounterpartyName: %v", transferKey.IntegerID(), fixer.transfer.Creator().ContactName)
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
