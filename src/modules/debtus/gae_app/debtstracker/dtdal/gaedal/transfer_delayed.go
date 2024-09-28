package gaedal

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"reflect"
	"sync"
	"time"
)

func (TransferDalGae) DelayUpdateTransfersWithCounterparty(ctx context.Context, creatorCounterpartyID, counterpartyCounterpartyID string) (err error) {
	logus.Debugf(ctx, "DelayUpdateTransfersWithCounterparty(creatorCounterpartyID=%s, counterpartyCounterpartyID=%s)", creatorCounterpartyID, counterpartyCounterpartyID)
	if creatorCounterpartyID == "" {
		return errors.New("creatorCounterpartyID == 0")
	}
	if counterpartyCounterpartyID == "" {
		return errors.New("counterpartyCounterpartyID == 0")
	}
	if err := delayerUpdateTransfersWithCounterparty.EnqueueWork(ctx, delaying.With(const4debtus.QueueTransfers, DELAY_UPDATE_TRANSFERS_WITH_COUNTERPARTY, 0), creatorCounterpartyID, counterpartyCounterpartyID); err != nil {
		return err
	}
	return nil
}

const (
	DELAY_UPDATE_TRANSFERS_WITH_COUNTERPARTY  = "update-api4transfers-with-counterparty"
	DELAY_UPDATE_1_TRANSFER_WITH_COUNTERPARTY = "update-1-transfer-with-counterparty"
)

func delayedUpdateTransfersWithCounterparty(ctx context.Context, creatorCounterpartyID, counterpartyCounterpartyID int64) (err error) {
	logus.Infof(ctx, "delayerUpdateTransfersWithCounterparty(creatorCounterpartyID=%d, counterpartyCounterpartyID=%d)", creatorCounterpartyID, counterpartyCounterpartyID)
	if creatorCounterpartyID == 0 {
		logus.Errorf(ctx, "creatorCounterpartyID == 0")
		return nil
	}
	if counterpartyCounterpartyID == 0 {
		logus.Errorf(ctx, "counterpartyCounterpartyID == 0")
		return nil
	}

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	query := dal.From(models4debtus.TransfersCollection).
		WhereField("BothCounterpartyIDs", dal.Equal, creatorCounterpartyID).WhereField("BothCounterpartyIDs", dal.Equal, 0).
		OrderBy(dal.DescendingField("DtCreated")).
		SelectKeysOnly(reflect.Int)

	var reader dal.Reader
	if reader, err = db.QueryReader(ctx, query); err != nil {
		return err
	}
	if transferIDs, err := dal.SelectAllIDs[int](reader, dal.WithLimit(query.Limit())); err != nil {
		return fmt.Errorf("failed to load api4transfers: %w", err)
	} else if len(transferIDs) > 0 {
		logus.Infof(ctx, "Loaded %d transfer IDs", len(transferIDs))
		delayDuration := 10 * time.Microsecond
		for _, transferID := range transferIDs {
			if err := delayerUpdateTransferWithCounterparty.EnqueueWork(ctx, delaying.With(const4debtus.QueueTransfers, DELAY_UPDATE_1_TRANSFER_WITH_COUNTERPARTY, delayDuration), transferID, counterpartyCounterpartyID); err != nil {
				return fmt.Errorf("failed to create task for transfer id=%d: %w", transferID, err)
			}
			delayDuration += 10 * time.Microsecond
		}
	} else {
		query := dal.From(models4debtus.TransfersCollection).
			WhereField("BothCounterpartyIDs", dal.Equal, creatorCounterpartyID).WhereField("BothCounterpartyIDs", dal.Equal, counterpartyCounterpartyID).
			Limit(1).
			SelectKeysOnly(reflect.Int)
		var reader dal.Reader
		if reader, err = db.QueryReader(ctx, query); err != nil {
			return err
		}
		var transferIDs []int
		if transferIDs, err = dal.SelectAllIDs[int](reader, dal.WithLimit(query.Limit())); err != nil {
			return fmt.Errorf("failed to load api4transfers by 2 counterparty IDs: %w", err)
		}
		if len(transferIDs) > 0 {
			logus.Infof(ctx, "No api4transfers found to update counterparty details")
		} else {
			logus.Warningf(ctx, "No api4transfers found to update counterparty details")
		}
	}
	return nil
}

func delayedUpdateTransferWithCounterparty(ctx context.Context, spaceID, transferID string, counterpartyCounterpartyID string) (err error) {
	logus.Debugf(ctx, "delayerUpdateTransferWithCounterparty(transferID=%s, counterpartyCounterpartyID=%s)", transferID, counterpartyCounterpartyID)
	if transferID == "" {
		logus.Errorf(ctx, "transferID == 0")
		return nil
	}
	if counterpartyCounterpartyID == "" {
		logus.Errorf(ctx, "counterpartyCounterpartyID == 0")
		return nil
	}

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return err
	}

	counterpartyCounterpartyContact := dal4contactus.NewContactEntry(spaceID, counterpartyCounterpartyID)
	if err = dal4contactus.GetContact(ctx, db, counterpartyCounterpartyContact); err != nil {
		logus.Errorf(ctx, err.Error())
		if dal.IsNotFound(err) {
			return nil
		}
		return err
	}

	counterpartyCounterpartyDebtusContact := models4debtus.NewDebtusSpaceContactEntry(spaceID, counterpartyCounterpartyID, nil)
	if err = facade4debtus.GetDebtusSpaceContact(ctx, db, counterpartyCounterpartyDebtusContact); err != nil {
		logus.Errorf(ctx, err.Error())
		if dal.IsNotFound(err) {
			return nil
		}
		return err
	}

	logus.Debugf(ctx, "counterpartyCounterpartyDebtusContact: %v", counterpartyCounterpartyDebtusContact)

	counterpartyUser := dbo4userus.NewUserEntry(counterpartyCounterpartyContact.Data.UserID)

	if err = dal4userus.GetUser(ctx, db, counterpartyUser); err != nil {
		logus.Errorf(ctx, err.Error())
		if dal.IsNotFound(err) {
			return nil
		}
		return err
	}

	logus.Debugf(ctx, "counterpartyUser: %v", *counterpartyUser.Data)

	if err := facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		transfer, err := facade4debtus.Transfers.GetTransferByID(ctx, tx, transferID)
		if err != nil {
			return err
		}
		changed := false

		// TODO: allow to pass creator counterparty as well. Match by userID

		logus.Debugf(ctx, "transfer.From() before: %v", transfer.Data.From())
		logus.Debugf(ctx, "transfer.To() before: %v", transfer.Data.To())

		// Update transfer creator
		{
			transferCreator := transfer.Data.Creator()
			logus.Debugf(ctx, "transferCreator before: %v", transferCreator)
			if transferCreator.ContactID == "" {
				transferCreator.ContactID = counterpartyCounterpartyDebtusContact.ID
				changed = true
			} else if transferCreator.ContactID != counterpartyCounterpartyDebtusContact.ID {
				err = fmt.Errorf("transferCounterparty.ContactID != counterpartyCounterpartyDebtusContact.ContactID: %s != %s", transferCreator.ContactID, counterpartyCounterpartyDebtusContact.ID)
				return err
			} else {
				logus.Debugf(ctx, "transferCounterparty.ContactID == counterpartyCounterpartyDebtusContact.ContactID: %s", transferCreator.ContactID)
			}
			if transferCreator.ContactName == "" || transferCreator.ContactName != counterpartyCounterpartyDebtusContact.Data.FullName() {
				transferCreator.ContactName = counterpartyCounterpartyDebtusContact.Data.FullName()
				changed = true
			}
			logus.Debugf(ctx, "transferCreator after: %v", transferCreator)
			logus.Debugf(ctx, "transfer.Creator() after: %v", transfer.Data.Creator())
		}

		// Update transfer counterparty
		{
			transferCounterparty := transfer.Data.Counterparty()
			logus.Debugf(ctx, "transferCounterparty before: %v", transferCounterparty)
			if transferCounterparty.UserID == "" {
				transferCounterparty.UserID = counterpartyCounterpartyContact.Data.UserID
				changed = true
			} else if transferCounterparty.UserID != counterpartyCounterpartyContact.Data.UserID {
				err = fmt.Errorf("transferCounterparty.UserID != counterpartyCounterpartyDebtusContact.UserID: %s != %s", transferCounterparty.UserID, counterpartyCounterpartyContact.Data.UserID)
				return err
			} else {
				logus.Debugf(ctx, "transferCounterparty.UserID == counterpartyCounterpartyDebtusContact.UserID: %s", transferCounterparty.UserID)
			}
			if transferCounterparty.UserName == "" || transferCounterparty.UserName != counterpartyUser.Data.Names.GetFullName() {
				transferCounterparty.UserName = counterpartyUser.Data.Names.GetFullName()
				changed = true
			}
			logus.Debugf(ctx, "transferCounterparty after: %v", transferCounterparty)
			logus.Debugf(ctx, "transfer.DebtusSpaceContactEntry() after: %v", transfer.Data.Counterparty())
		}
		logus.Debugf(ctx, "transfer.From() after: %v", transfer.Data.From())
		logus.Debugf(ctx, "transfer.To() after: %v", transfer.Data.To())

		if changed {
			if err = facade4debtus.Transfers.SaveTransfer(ctx, tx, transfer); err != nil {
				return err
			}
			if !transfer.Data.DtDueOn.IsZero() {
				counterpartyDebtusSpace := models4debtus.NewDebtusSpaceEntry(spaceID)
				if err = models4debtus.GetDebtusSpace(ctx, tx, counterpartyDebtusSpace); err != nil {
					return err
				}

				if !counterpartyDebtusSpace.Data.HasDueTransfers {
					if err = facade4debtus.DelayUpdateHasDueTransfers(ctx, counterpartyCounterpartyContact.Data.UserID, spaceID); err != nil {
						return err
					}
				}
			}
			logus.Infof(ctx, "TransferEntry saved to datastore")
			return nil
		} else {
			logus.Infof(ctx, "No changes for the transfer")
		}
		return nil
	}, nil); err != nil {
		panic(fmt.Sprintf("failed to update transfer (%s): %v", transferID, err.Error()))
	} else {
		logus.Infof(ctx, "Transaction successfully completed")
	}
	return nil
}

const (
	UPDATE_TRANSFERS_WITH_CREATOR_NAME = "update-api4transfers-with-creator-name"
)

func DelayUpdateTransfersWithCreatorName(ctx context.Context, userID string) error {
	return delayerUpdateTransfersWithCreatorName.EnqueueWork(ctx, delaying.With(const4debtus.QueueTransfers, UPDATE_TRANSFERS_WITH_CREATOR_NAME, 0), userID)
}

func delayedUpdateTransfersWithCreatorName(ctx context.Context, userID string) (err error) {
	logus.Debugf(ctx, "delayedUpdateTransfersWithCreatorName(userID=%s)", userID)

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return err
	}
	user := dbo4userus.NewUserEntry(userID)
	if err = dal4userus.GetUser(ctx, db, user); err != nil {
		logus.Errorf(ctx, err.Error())
		if dal.IsNotFound(err) {
			err = nil
		}
		return err
	}

	userName := user.Data.Names.GetFullName()

	query := dal.From(models4debtus.TransfersCollection).
		WhereField("BothUserIDs", dal.Equal, userID).
		SelectInto(models4debtus.NewTransferRecord)

	var reader dal.Reader
	reader, err = db.QueryReader(ctx, query)

	var wg sync.WaitGroup
	defer wg.Wait()
	for {
		transferRecord, err := reader.Next()
		if err != nil {
			return err
		}
		transfer := models4debtus.TransferFromRecord(transferRecord)
		if err != nil {
			if errors.Is(err, dal.ErrNoMoreRecords) {
				return nil
			}
			logus.Errorf(ctx, err.Error())
			return err
		}
		wg.Add(1)
		go func(transferID string) {
			defer wg.Done()
			err := facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
				transfer, err := facade4debtus.Transfers.GetTransferByID(ctx, tx, transferID)
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
					logus.Infof(ctx, "TransferEntry() creator is not a counterparty")
				}
				if changed {
					if err = facade4debtus.Transfers.SaveTransfer(ctx, tx, transfer); err != nil {
						return err
					}
				}
				return err
			})
			if err != nil {
				logus.Errorf(ctx, err.Error())
			}
		}(transfer.ID)
	}
}
