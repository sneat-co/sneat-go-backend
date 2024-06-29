package gaedal

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"reflect"
	"sync"
	"time"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func (TransferDalGae) DelayUpdateTransfersWithCounterparty(c context.Context, creatorCounterpartyID, counterpartyCounterpartyID string) (err error) {
	logus.Debugf(c, "DelayUpdateTransfersWithCounterparty(creatorCounterpartyID=%s, counterpartyCounterpartyID=%s)", creatorCounterpartyID, counterpartyCounterpartyID)
	if creatorCounterpartyID == "" {
		return errors.New("creatorCounterpartyID == 0")
	}
	if counterpartyCounterpartyID == "" {
		return errors.New("counterpartyCounterpartyID == 0")
	}
	if err := delayUpdateTransfersWithCounterparty.EnqueueWork(c, delaying.With(common.QUEUE_TRANSFERS, DELAY_UPDATE_TRANSFERS_WITH_COUNTERPARTY, 0), creatorCounterpartyID, counterpartyCounterpartyID); err != nil {
		return err
	}
	return nil
}

const (
	DELAY_UPDATE_TRANSFERS_WITH_COUNTERPARTY  = "update-transfers-with-counterparty"
	DELAY_UPDATE_1_TRANSFER_WITH_COUNTERPARTY = "update-1-transfer-with-counterparty"
)

func delayedUpdateTransfersWithCounterparty(c context.Context, creatorCounterpartyID, counterpartyCounterpartyID int64) (err error) {
	logus.Infof(c, "delayUpdateTransfersWithCounterparty(creatorCounterpartyID=%d, counterpartyCounterpartyID=%d)", creatorCounterpartyID, counterpartyCounterpartyID)
	if creatorCounterpartyID == 0 {
		logus.Errorf(c, "creatorCounterpartyID == 0")
		return nil
	}
	if counterpartyCounterpartyID == 0 {
		logus.Errorf(c, "counterpartyCounterpartyID == 0")
		return nil
	}

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	query := dal.From(models.TransfersCollection).
		WhereField("BothCounterpartyIDs", dal.Equal, creatorCounterpartyID).WhereField("BothCounterpartyIDs", dal.Equal, 0).
		OrderBy(dal.DescendingField("DtCreated")).
		SelectKeysOnly(reflect.Int)

	var reader dal.Reader
	if reader, err = db.QueryReader(c, query); err != nil {
		return err
	}
	if transferIDs, err := dal.SelectAllIDs[int](reader, dal.WithLimit(query.Limit())); err != nil {
		return fmt.Errorf("failed to load transfers: %w", err)
	} else if len(transferIDs) > 0 {
		logus.Infof(c, "Loaded %d transfer IDs", len(transferIDs))
		delayDuration := 10 * time.Microsecond
		for _, transferID := range transferIDs {
			if err := delayUpdateTransferWithCounterparty.EnqueueWork(c, delaying.With(common.QUEUE_TRANSFERS, DELAY_UPDATE_1_TRANSFER_WITH_COUNTERPARTY, delayDuration), transferID, counterpartyCounterpartyID); err != nil {
				return fmt.Errorf("failed to create task for transfer id=%d: %w", transferID, err)
			}
			delayDuration += 10 * time.Microsecond
		}
	} else {
		query := dal.From(models.TransfersCollection).
			WhereField("BothCounterpartyIDs", dal.Equal, creatorCounterpartyID).WhereField("BothCounterpartyIDs", dal.Equal, counterpartyCounterpartyID).
			Limit(1).
			SelectKeysOnly(reflect.Int)
		var reader dal.Reader
		if reader, err = db.QueryReader(c, query); err != nil {
			return err
		}
		var transferIDs []int
		if transferIDs, err = dal.SelectAllIDs[int](reader, dal.WithLimit(query.Limit())); err != nil {
			return fmt.Errorf("failed to load transfers by 2 counterparty IDs: %w", err)
		}
		if len(transferIDs) > 0 {
			logus.Infof(c, "No transfers found to update counterparty details")
		} else {
			logus.Warningf(c, "No transfers found to update counterparty details")
		}
	}
	return nil
}

func delayedUpdateTransferWithCounterparty(c context.Context, transferID string, counterpartyCounterpartyID string) (err error) {
	logus.Debugf(c, "delayUpdateTransferWithCounterparty(transferID=%s, counterpartyCounterpartyID=%s)", transferID, counterpartyCounterpartyID)
	if transferID == "" {
		logus.Errorf(c, "transferID == 0")
		return nil
	}
	if counterpartyCounterpartyID == "" {
		logus.Errorf(c, "counterpartyCounterpartyID == 0")
		return nil
	}

	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return err
	}

	counterpartyCounterparty, err := facade.GetContactByID(c, db, counterpartyCounterpartyID)
	if err != nil {
		logus.Errorf(c, err.Error())
		if dal.IsNotFound(err) {
			return nil
		}
		return err
	}

	logus.Debugf(c, "counterpartyCounterparty: %v", counterpartyCounterparty)

	counterpartyUser, err := facade.User.GetUserByID(c, db, counterpartyCounterparty.Data.UserID)
	if err != nil {
		logus.Errorf(c, err.Error())
		if dal.IsNotFound(err) {
			return nil
		}
		return err
	}

	logus.Debugf(c, "counterpartyUser: %v", *counterpartyUser.Data)

	if err := db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		transfer, err := facade.Transfers.GetTransferByID(tc, tx, transferID)
		if err != nil {
			return err
		}
		changed := false

		// TODO: allow to pass creator counterparty as well. Match by userID

		logus.Debugf(c, "transfer.From() before: %v", transfer.Data.From())
		logus.Debugf(c, "transfer.To() before: %v", transfer.Data.To())

		// Update transfer creator
		{
			transferCreator := transfer.Data.Creator()
			logus.Debugf(c, "transferCreator before: %v", transferCreator)
			if transferCreator.ContactID == "" {
				transferCreator.ContactID = counterpartyCounterparty.ID
				changed = true
			} else if transferCreator.ContactID != counterpartyCounterparty.ID {
				err = fmt.Errorf("transferCounterparty.ContactID != counterpartyCounterparty.ID: %s != %s", transferCreator.ContactID, counterpartyCounterparty.ID)
				return err
			} else {
				logus.Debugf(c, "transferCounterparty.ContactID == counterpartyCounterparty.ID: %s", transferCreator.ContactID)
			}
			if transferCreator.ContactName == "" || transferCreator.ContactName != counterpartyCounterparty.Data.FullName() {
				transferCreator.ContactName = counterpartyCounterparty.Data.FullName()
				changed = true
			}
			logus.Debugf(c, "transferCreator after: %v", transferCreator)
			logus.Debugf(c, "transfer.Creator() after: %v", transfer.Data.Creator())
		}

		// Update transfer counterparty
		{
			transferCounterparty := transfer.Data.Counterparty()
			logus.Debugf(c, "transferCounterparty before: %v", transferCounterparty)
			if transferCounterparty.UserID == "" {
				transferCounterparty.UserID = counterpartyCounterparty.Data.UserID
				changed = true
			} else if transferCounterparty.UserID != counterpartyCounterparty.Data.UserID {
				err = fmt.Errorf("transferCounterparty.UserID != counterpartyCounterparty.UserID: %s != %s", transferCounterparty.UserID, counterpartyCounterparty.Data.UserID)
				return err
			} else {
				logus.Debugf(c, "transferCounterparty.UserID == counterpartyCounterparty.UserID: %s", transferCounterparty.UserID)
			}
			if transferCounterparty.UserName == "" || transferCounterparty.UserName != counterpartyUser.Data.FullName() {
				transferCounterparty.UserName = counterpartyUser.Data.FullName()
				changed = true
			}
			logus.Debugf(c, "transferCounterparty after: %v", transferCounterparty)
			logus.Debugf(c, "transfer.ContactEntry() after: %v", transfer.Data.Counterparty())
		}
		logus.Debugf(c, "transfer.From() after: %v", transfer.Data.From())
		logus.Debugf(c, "transfer.To() after: %v", transfer.Data.To())

		if changed {
			if err = facade.Transfers.SaveTransfer(tc, tx, transfer); err != nil {
				return err
			}
			if !transfer.Data.DtDueOn.IsZero() {
				var counterpartyUser models.AppUser
				if counterpartyUser, err = facade.User.GetUserByID(c, tx, counterpartyCounterparty.Data.UserID); err != nil {
					return err
				}

				if !counterpartyUser.Data.HasDueTransfers {
					if err = dtdal.User.DelayUpdateUserHasDueTransfers(tc, counterpartyCounterparty.Data.UserID); err != nil {
						return err
					}
				}
			}
			logus.Infof(c, "TransferEntry saved to datastore")
			return nil
		} else {
			logus.Infof(c, "No changes for the transfer")
		}
		return nil
	}, nil); err != nil {
		panic(fmt.Sprintf("failed to update transfer (%s): %v", transferID, err.Error()))
	} else {
		logus.Infof(c, "Transaction successfully completed")
	}
	return nil
}

const (
	UPDATE_TRANSFERS_WITH_CREATOR_NAME = "update-transfers-with-creator-name"
)

func DelayUpdateTransfersWithCreatorName(c context.Context, userID string) error {
	return delayUpdateTransfersWithCreatorName.EnqueueWork(c, delaying.With(common.QUEUE_TRANSFERS, UPDATE_TRANSFERS_WITH_CREATOR_NAME, 0), userID)
}

func delayedUpdateTransfersWithCreatorName(c context.Context, userID string) (err error) {
	logus.Debugf(c, "delayedUpdateTransfersWithCreatorName(userID=%s)", userID)

	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return err
	}
	user, err := facade.User.GetUserByID(c, db, userID)
	if err != nil {
		logus.Errorf(c, err.Error())
		if dal.IsNotFound(err) {
			err = nil
		}
		return err
	}

	userName := user.Data.FullName()

	query := dal.From(models.TransfersCollection).
		WhereField("BothUserIDs", dal.Equal, userID).
		SelectInto(models.NewTransferRecord)

	var reader dal.Reader
	reader, err = db.QueryReader(c, query)

	var wg sync.WaitGroup
	defer wg.Wait()
	for {
		transferRecord, err := reader.Next()
		if err != nil {
			return err
		}
		transfer := models.TransferFromRecord(transferRecord)
		if err != nil {
			if err == dal.ErrNoMoreRecords {
				return nil
			}
			logus.Errorf(c, err.Error())
			return err
		}
		wg.Add(1)
		go func(transferID string) {
			defer wg.Done()
			err := db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
				transfer, err := facade.Transfers.GetTransferByID(c, tx, transferID)
				if err != nil {
					return err
				}
				changed := false
				switch userID {
				case transfer.Data.From().UserID:
					if from := transfer.Data.From(); from.UserName != userName {
						from.UserName = userName
						changed = true
					}
				case transfer.Data.To().UserID:
					if to := transfer.Data.To(); to.UserName != userName {
						to.UserName = userName
						changed = true
					}
				default:
					logus.Infof(c, "TransferEntry() creator is not a counterparty")
				}
				if changed {
					if err = facade.Transfers.SaveTransfer(c, tx, transfer); err != nil {
						return err
					}
				}
				return err
			})
			if err != nil {
				logus.Errorf(c, err.Error())
			}
		}(transfer.ID)
	}
}
